#include "docker_api.h"
#include <iostream>
#include <curl/curl.h>

// using namespace for nlohmann::json is common and acceptable here
using json = nlohmann::json;

// --- Helper Functions ---

// libcurl write callback function
// It appends the received data to a std::string
size_t WriteCallback(void *contents, size_t size, size_t nmemb, void *userp)
{
    ((std::string *)userp)->append((char *)contents, size * nmemb);
    return size * nmemb;
}

// --- DockerAPI Class Implementation ---

json DockerAPI::getContainerInfo()
{
    CURL *curl;
    CURLcode res;
    std::string readBuffer;
    json containersJson;

    curl = curl_easy_init();
    if (curl)
    {
        std::string url = "http://localhost/v1.41/containers/json?filters={\"status\":[\"running\"]}";

        curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
        curl_easy_setopt(curl, CURLOPT_UNIX_SOCKET_PATH, "/run/docker.sock");

        // This is a crucial line for the JSON response to be stored
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &readBuffer);

        res = curl_easy_perform(curl);

        if (res != CURLE_OK)
        {
            std::cerr << "curl_easy_perform() failed: " << curl_easy_strerror(res) << std::endl;
        }
        else
        {
            try
            {
                containersJson = json::parse(readBuffer);
            }
            catch (const json::parse_error &e)
            {
                std::cerr << "JSON parsing failed: " << e.what() << std::endl;
            }
        }

        // Always cleanup
        curl_easy_cleanup(curl);
    }

    return containersJson;
}

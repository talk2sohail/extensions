#ifndef DOCKER_API_H
#define DOCKER_API_H

#include <string>
#include <curl/curl.h>
#include <nlohmann/json.hpp>

class DockerAPI
{
public:
    nlohmann::json getContainerInfo();

private:
    CURL *curl;
};

#endif // DOCKER_API_H
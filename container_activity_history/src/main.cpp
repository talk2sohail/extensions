#include <iostream>
#include "docker_api.h"

int main()
{
    std::cout << "Docker API C++ Application" << std::endl;

    DockerAPI dockerApi;

    try
    {
        nlohmann::json containerList = dockerApi.getContainerInfo();

        if (containerList.is_array())
        {
            if (!containerList.empty())
            {
                std::cout << "Running Containers:" << std::endl;
                for (const auto &container : containerList)
                {
                    std::string id = container["Id"];
                    std::string name = container["Names"][0];
                    std::string image = container["Image"];

                    std::cout << "  - ID: " << id.substr(0, 12)
                              << ", Name: " << name.substr(1)
                              << ", Image: " << image << std::endl;
                }
            }
            else
            {
                std::cout << "No running containers found." << std::endl;
            }
        }
        else
        {
            std::cout << "Failed to get container information. Is the Docker daemon running and accessible?" << std::endl;
        }
    }
    catch (const std::exception &e)
    {
        std::cerr << "An error occurred: " << e.what() << std::endl;
    }

    return 0;
}
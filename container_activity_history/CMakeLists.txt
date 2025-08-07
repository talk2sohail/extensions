cmake_minimum_required(VERSION 3.10)

# Set the project name
project(DockerCppApp)

# Specify the C++ standard
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED True)

# Include directories
include_directories(include)

# Find required packages
find_package(CURL REQUIRED)
find_package(nlohmann_json 3.2.0 REQUIRED)

# Add the executable
add_executable(DockerCppApp src/main.cpp src/docker_api.cpp)

# Link the libraries
target_link_libraries(DockerCppApp PRIVATE CURL::libcurl nlohmann_json::nlohmann_json)
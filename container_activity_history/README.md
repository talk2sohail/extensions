# Docker C++ API Application

This project is a simple C++ application that interacts with the Docker Engine API to retrieve information about containers.

## Project Structure

```
docker-cpp-app
├── src
│   ├── main.cpp          # Entry point of the application
│   └── docker_api.cpp    # Implementation of Docker API functions
├── include
│   └── docker_api.h      # Header file for Docker API functions
├── CMakeLists.txt        # CMake configuration file
└── README.md             # Project documentation
```

## Setup Instructions

### Prerequisites

Make sure you have CMake installed on your system.

### Install Libraries

1. **Install libcurl:**
   - On Ubuntu:
     ```
     sudo apt-get install libcurl4-openssl-dev
     ```
   - On macOS:
     ```
     brew install curl
     ```

2. **Install nlohmann/json:**
   - You can add it as a dependency in your `CMakeLists.txt` by including the following line:
     ```
     find_package(nlohmann_json 3.2.0 REQUIRED)
     ```

### Update CMakeLists.txt

Make sure to link the libraries in your `CMakeLists.txt`:
```cmake
target_link_libraries(your_target_name PRIVATE curl nlohmann_json::nlohmann_json)
```

### Build the Project

1. Create a build directory:
   ```
   mkdir build && cd build
   ```

2. Run CMake:
   ```
   cmake ..
   ```

3. Compile the project:
   ```
   make
   ```

This will set up the build system for your C++ application that interacts with the Docker API.
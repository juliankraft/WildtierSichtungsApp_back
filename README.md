## WildTierSichtungsApp - Backend

**Author:**         Julian Kraft   
**Institution:**    Zurich University of Applied Sciences (ZHAW)
**Program:**        BSc Natural Resource Sciences
**Course:**         Angewandte Geoinformatik
**Project:**        Semester Project  
**Date:**           2024-12-18

### Abstract

The project aimed to evaluate the accessibility of modern technology for creating a web application 
that records geotagged data. Using modern web frameworks, a prototype was developed to document wildlife sightings. 
This application, accessible online and installable on smartphones across major operating systems, 
exemplifies how a straightforward and purpose-driven tool can be built with moderate programming knowledge, 
guided support, and persistence. The project's workflow encompassed backend development in Go, 
frontend design with the Ionic framework, and database integration using MariaDB. Additionally, 
the system incorporates a public data visualization tool and direct database access. 
While the project highlights the potential of current technologies and AI-powered tools like GitHub Copilot, 
it underscores challenges in managing interconnected components. The result is a versatile platform, 
demonstrating both the promise and complexity of custom web application development for geospatial data collection.

### Repository Content

This reository contains the backend of the WildTierSichtungsApp, the compiled version of the frontend and a static webpage.

### Repository Structure

/app                        # compiled version of the frontend
/db_setup                   # scripts to setup the database and access data
/static                     # static webpage accessible from web
go.mod                      # Go module file
go.sum                      # Go module checksum file
handlers_dataentry.go       # Go file with handlers for data entry
handlers_login.go           # Go file with handlers for login
mandlers_table.go           # Go file with handlers data view
main.go                     # Go file with main function
middleware.go               # Go file with middleware functions
server.sh                   # script to run the server
README.md                   # Project documentation
LICENSE.md                  # License information

### Additional Resources:

- Frontend: [WildTierSichtungsApp-Frontend](https://github.com/juliankraft/WildtierSichtungsApp_front)
- Documentation: [WildTierSichtungsApp-Documentation](https://github.com/juliankraft/WildtierSichtungsApp_documentation)

### License

This project is licensed under the Creative Commons Zero v1.0 Universal License - see the [LICENSE.md](LICENSE.md) file for details.
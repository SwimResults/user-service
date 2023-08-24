# user-service

microservice with REST API to fetch additional data to users, store followers and user settings

## Endpoints

    [GIN-debug] GET    /user                     --> github.com/swimresults/user-service/controller.getUsers (3 handlers)
    [GIN-debug] GET    /user/:id                 --> github.com/swimresults/user-service/controller.getUser (3 handlers)   
    [GIN-debug] DELETE /user/:id                 --> github.com/swimresults/user-service/controller.removeUser (3 handlers)
    [GIN-debug] POST   /user                     --> github.com/swimresults/user-service/controller.addUser (3 handlers)   
    [GIN-debug] PUT    /user                     --> github.com/swimresults/user-service/controller.updateUser (3 handlers)
    [GIN-debug] GET    /actuator                 --> github.com/swimresults/user-service/controller.actuator (3 handlers)
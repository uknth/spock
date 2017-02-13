# Spock - Request Router

### Goals
- Routing requests based on dynamic states
- Volume shouldn’t affect the performance of the system
- Shouldn’t add additional latency
- Should follow the guidelines set by RFC 7230-37

### Specification
- **Listener**   
Listener actively listens to incoming request on network interface. Once a requesthits, it checks forwards it to loaded middlewares (aka filters) and then passes it toHandlers to be handled.
- **Middlewares**   
Middlewares augment the request/response object. For instance, it can be used toasync log a req/res or can be used to update headers etc
- **Handler**   
Handler handles the request received. It is an entry point of application and isexecuted right after some `pre` middlewares (**Note**: *Post middlewares are executed after handler is executed*). It is equivalent to servlets in java. Any requestthat is received is then further sent to services/engine to be handled for businesslogic. Any network related changes or object creation is managed here
- **Engine**   
Engine comprises of multiple subdivisions. It is a wrapper around all internalpackages which do certain set of workIts responsibility is to gather request from handlers (HandleReq() method) and dorouting based on the routing strategy and virtual server configuration   
   
 - **Config**   
Config exposes an interface of `ResourceLoader` which is responsible ofloading the configuration from a source.There are multiple instances of ResourceLoader possible, which can loadconfigurations from external data source or from state machines likeZookeeper.It gets triggered by Engine Loader and is made available in Virtual Server.
 - **HA (high availability) Manager**   
HA Manager is responsible for High Availability of downstream. There arescenarios where a decision should be made on the performance ofdownstream and a new instance of VM or Container needs to be initializedto rebalance the load.This acts as an interface for that, and conditional rules responsible forrebalance is loaded per Virtual Server basis. When a rebalance scenariooccurs for a particular Virtual Server, dependent action is triggered for it.
 - **HealthChecker**   
HealthChecker interface has multiple implementations possible and aparticular health checker is loaded for a single instance of Virtual Server.The rules of Health Check are defined in the configuration and an instanceof it is triggered per Downstream host.
 - **Router**   
Router is responsible for actual routing of the request. A custom strategyfor routing is loaded per Virtual Server and a request is then routed to acertain destination based on it. The Router doesn’t make a call todownstream, instead it updates the necessary request parameter withrouting information
 - **VServer (Virtual Server)**
Virtual Server is an object which contains configuration, HA, Health Checkerand Router object for a particular domain. This object is maintained perdomain and will be used to decide which downstream to call or what is thestatus of the health check or to check if we need to rebalance thedownstream etc.This is closest thing to a `Core` and has objects which will be used toprocess the request
 - **Sync**   
Sync is responsible for maintaining states across multiple Spock nodes. It isresponsible for Leader Election, and leader will then make decision onfailover or HA.
 - **IPC**   
IPC is responsible for communication between different processes runningon same machine. It facilitates communications between all externalcommand line tools such as ones responsible for shutting down, etc.It also exposes API for statistics etc
 - **Dialer**   
Actual request to downstream is managed by the dialer package. The Engineforms a curated Request object based on the parameters in VServers and sends itto Dialer. Dialer will make the actual `network` call to downstream and will `fetch` theresponse back

### Functionality
Request is received by the listener which then propagates it through series ofmiddlewares loaded.Middleware performs augmentation on the request object and sends it further to therequest handlerRequest handler then sends the request to Engine.HandleReq() which picks the route tochoose based on available parameters in the matched VServer. If there is no match for aVServer present, it should throw a 404.Engine creates a new Request object to be sent to the downstream and appends it withnecessary informations based on info from Router & Health Check.Engine sends the request object to Dialer which makes the downstream call, gets theresponse and returns it back.

# kube8-operator

## About

The Operator Pattern in Kubernetes is a method for extending the functionality of Kubernetes by defining custom resources and controllers. It allows us to automate the management of complex applications and services by encapsulating operational knowledge into a Kubernetes-native application. With this extension, we can manage and automate the lifecycle of specific applications or services.

The Operator typically consists of two main components:
- Custom Resource Definition (CRD): Defines a new custom resource that extends the Kubernetes API. It represents the application or service we want to manage using the controller. In this case, the custom resource is the [collector resource](internal/operator/crd.yaml). This includes the structure and behavior of the custom resource. It also defines the API group, version, and kind for the custom resource. We then outline the desired fields, their types, validation rules, and any additional behaviors we want to associate with the resource.
- Controller: This is the logic that runs within the Operator. It watches for changes to the custom collector resource and takes actions accordingly to reconcile the desired state with the actual state. The controller is responsible for grabbing and deploying the helm chart associated with the collector resource.

### Collector Resource
The collector resource is a custom resource that represents the collector. It defines the desired state of the collector and is used by the controller to reconcile the actual state with the desired state. The collector resource is defined in the [collector CRD](internal/operator/crd.yaml). The Collector resource is used to fill in the details of the collector configuration and the helm deployment request so that the collector resources (deployment, secrets, service, serviceMonitor, etc. from the chart) can be created in the cluster.

The Collector resource is structured with several properties within its schema:
- **apiVersion** and **kind** as standard Kubernetes resource properties.
- **metadata** for managing metadata of the resource.
- **spec** defining the desired state of the Collector with sub-properties like **collector** (collector name, version, and configuration), **tenant**, and **cluster**.
- **status** representing the observed state of the Collector with nested properties like conditions detailing status conditions.

### Controller Initialization
The NewController function initializes a controller instance that manages interactions with the Kubernetes API and handles events related to changes in the Collector resource. It also sets up the informer factory to receive notifications about changes in the collector resource.

Components Initialization:
- **Clients Creation**: Initializes clients for Kubernetes API operations.
- **Informer Factory**: Creates an informer factory to receive notifications about collector changes.
- **Event Broadcaster**: Sets up an event broadcaster to record events associated with the controller.
- **Event Recorder**: Creates an event recorder to log events within the Kubernetes system.

Controller Setup:
- **Work Queue**: Establishes a work queue for event handling.
- **Controller Object**: Initializes the controller object with various essential components such as Kubernetes clients, informers, event recorder, and observer.

Event Handling:
- **Event Handlers**: Defines event handler functions for different types of events (Add, Update, Delete) within the informer.

Execution:
- **Informer Factory Start**: Initiates the informer factory to start receiving and processing events.

### Processing Resource Creation/Update/Deletion by Controller

1. Resource Creation Process:
    - **Resource Creation**: A new instance of the Collector resource is created in the Kubernetes cluster that the operator is deployed to.
    - **API Server Handling**: The Kubernetes API server receives and validates the resource creation request against the defined CRD schema to ensure compliance and correctness.
2. Controller's Reaction to Resource Creation:
    - **Informer Watches for Changes**: The Controller's informer, initialized by the NewController function, constantly watches for changes related to Collector resources.
    - **Event Handling**: Upon the creation of an Collector resource, the informer detects this change and triggers the corresponding event handler in the Controller (addFunc, updateFunc, deleteFunc). The event handler then calls the appropriate function to reconcile the desired state with the actual state.
3. Controller's Reconciliation Process:
    - **Helm Client**: A helm client is created to send requests to the cluster.
    - **Helm Chart Construction**: The CreateOrUpdateCollector method in the [reconciler](internal/operator/reconciler.go) file decodes the Base64 encoded collector configuration and constructs a helm request based on the Collector resource's configuration.
    - **Helm deployment Creation/Update/Deletion**: The helm request is sent to the cluster to create/update/delete the collector resources (deployment, secrets, service, serviceMonitor, etc.) in the tenant's namespace.

### Managing Custom Operator API Code Generation

- **pkg Directory**: Contains all API-related code for the custom operators. Generated clientset, informer, listers, Collector register schema, type definitions, and generated.deepcopy.go file. The generated api code is essential for custom operators to communicate to the kubernetes API server, utilize the CRD types, includes the informer and listers that monitor and track changes to custom resources, and register the custom resource with the scheme (a lot more to unpack here).
- **hack Directory**: Contains the code generation scripts and boilerplate code for the custom operators. The code generation scripts are used to generate the API code in the pkg directory. This script is responsible for updating generated code if changes occur in the Custom Resource Definition (CRD).
    - **Note**: This should be used sparingly. Unless a change is made to the CRD, the script should not be run. If the script is run, it will overwrite any changes made to the generated code.

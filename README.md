Scalable Go Microservice with Hot Reloading, gRPC, RabbitMQ, Docker, MongoDB, Jaeger, Consul & Stripe

This project implements a robust microservice architecture designed for scalability, maintainability, and resilience. It leverages modern technologies to facilitate rapid development, efficient communication, and reliable data management.

Key Features:

Modern Go 1.22: Utilizes the latest Go language features for performance and maintainability. 
Hot Reloading (cosmtrek/air): Enables rapid development iterations without service restarts. 
gRPC Communication: Provides high-performance inter-service communication. 
RabbitMQ Message Broker: Ensures asynchronous message delivery and decoupling for fault tolerance. 
Docker & Docker Compose: Simplifies deployment, scaling, and environment consistency. 
MongoDB Storage: Offers a flexible NoSQL solution for document-oriented data. 
Jaeger Service Tracing: Provides comprehensive insights into service calls and performance metrics. 
HashiCorp Consul Service Discovery: Enables dynamic service registration and discovery. 
Stripe Payment Integration: (Optional) Streamlines secure payment processing. 

Getting Started:

Local
Docker Compose
For external services like RabbitMQ and JaggerUI, you can use docker compose to start them up.

cd ..
docker compose up --build
Start the services
cd order && air
cd payment && air
...
Start Stripe Server
Run the following command to start the stripe cli

stripe login
Then run the following command to listen for webhooks

stripe listen --forward-to localhost:8081/webhook
Where localhost:8081/webhook is the endpoint payment service HTTP server address.

Test card: 4242424242424242

RabbitMQ UI
http://localhost:15672/#/

Jaeger UI
Deployment
Build Docker images for each microservice and push them to a container registry. Deploy using Docker Compose or orchestration tools like Kubernetes.



Future implementation :
Implement error handling and retry mechanisms for service resilience. 
Define an authentication and authorization strategy for secure access. 
Establish monitoring and alerting for proactive issue identification. 
Utilize version control (Git) and CI/CD pipelines for streamlined development and deployment. 
Consider Infrastructure as Code (IaC) tools like Terraform for managing infrastructure configurations. Customization:

This architecture can be adapted to your specific needs. Refer to the code and configuration files for customization options.

Feel free to contribute or raise issues through pull requests and discussions.

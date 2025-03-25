# Poll-Voting-Website
This is a full-stack web application where users can vote on a daily "hot" question that presents two answer choices – simple, fun, or thought-provoking. The platform is designed to be lightweight, interactive, and engaging.

🛠️ Tech Stack
- Frontend: NextJS
- Backend: Golang (Gin)
- Containerization: Docker
- Deployment: AWS ECS (Fargate)

☁️ AWS Services Used
- Amazon ECS (Fargate): For hosting both the frontend and backend in containers, providing a fully serverless infrastructure.
- Amazon ElastiCache (Redis): Used to store real-time voting data with 24-hour TTL for performance and automatic reset.
- Amazon SNS (Simple Notification Service): Sends notifications for vote milestones and new trending questions.
- Amazon RDS (PostgreSQL): Stores user authentication, vote metadata, and application content.

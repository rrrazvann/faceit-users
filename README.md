# Users manager

## Description

This is a simple app for managing users.

## Getting Started

### Installing

1. Start docker containers.
2. Create database.
3. Set config.yml.
4. Build binaries by running `make`.
5. Run API from `./bin/api`.

### Improvements:

- save event subscribers data in a persistent storage;
- webhook system needs a retry system;
- make the microservice asynchronous: implement a job system so every request returns a job_id so they can poll to get the job's result. Or also, accept a webhook URL to get notified.

# Workflow Engine


# Introduction

This is an experiment to impelment a simple workflow engine and to rediscovery key components in such system. As a reference, we will compare our implementation to Prefect.

We choose Go for the implementation for the following reasons:
- We want to expose challenges caused by the requirement to specify types in such systems. Python does not have such burden.
- We want to perform some simulations and Go is a good fit for such tasks.


# Building Block 1: Persistence

Before diving into some complicated system design conversations, the simple fact that a workflow engine needs to support failure and recovery means execution states need to be persisted. The persistence handling is obviously a standalone components and ideally we should be able to configure the backend. To simplify our experiment, we'll just use PostgreSql for state persistence.


# Building Block 2: Task and Workflow


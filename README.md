# go-amazon-sqs-sample

The sample about active-mq, which use these libraries:
- [aws-sdk-go](https://github.com/aws/aws-sdk-go)
- [core-go/sqs](https://github.com/core-go/sqs) to wrap [aws-sdk-go](https://github.com/aws/aws-sdk-go)
    - Simplify the way to initialize the consumer, publisher by configurations
        - Props: when you want to change the parameter of consumer or publisher, you can change the config file, and restart Kubernetes POD, do not need to change source code and re-compile.
- [core-go/mq](https://github.com/core-go/mq) to implement this flow, which can be considered a low code tool for message queue consumer:

  ![A common flow to consume a message from a message queue](https://cdn-images-1.medium.com/max/800/1*Y4QUN6QnfmJgaKigcNHbQA.png)

# NGINX Example

We'll use AWS Ec2 to spin up a t2 micro and use the ec2-user data to do initial download of bruce.
Bruce will use the install.yml file within this directory to install nginx and configure it with a default vhost
which redirects http -> https & also updates the standard html under /var/www/html/index.hml with some reconfigured variables.

NOTE: templates within this example should match directory wise with what will be created on the server to make it easier to inspect.

I used some basic terraform to create a generic VPC with a single subnet (look at TF on how to do it or make one yourself):
This example should work fine on an intel nuc etc.  To understand how that is done take a look at the ec2-userdata.txt file in this directory.

# Step 1

First download the ec2-userdata.txt file locally this will be used to bootstrap your ec2 with bruce.

```
wget https://raw.githubusercontent.com/Nitecon/bruce/main/examples/nginx/ec2-userdata.txt
```

After you have configured ec2-userdata copy the text from below (after you have configured aws cli)

Adjust the ssh key name / image id / subnets / sec groups etc to match your environment and then start up that instance

```
aws ec2 run-instances --image-id ami-09d3b3274b6c5d4aa --count 1 \
--instance-type t2.micro --key-name mynewkey \
--security-group-ids sg-0934ce3940ab515dc
--subnet-id subnet-052cc33d8f0f0c960 --user-data file://ec2-userdata.txt
```

After your ec2 instance command create is completed you can now hook it up to a load balancer to view the output, or simply log onto the box and run:
```
curl http://localhost/
```
This should give you a redirect to https://localhost, so we know it's working correctly,
Next is to view the output of the vhost that was configured with:

```
curl --header 'Host: www.example.com' --insecure https://127.0.0.1
```

Since the ssl cert we created is self signed we pass --insecure to make sure it loads.

- Enjoy
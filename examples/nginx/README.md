# NGINX Example

We'll use AWS Ec2 to spin up a t2 micro and use the ec2-user data to do initial download of bruce.
Bruce will use the install.yml file within this directory to install nginx and configure it with a default vhost
which redirects http -> https & also updates the standard html under /var/www/html/index.hml with some reconfigured variables.

NOTE: templates within this example should match directory wise with what will be created on the server to make it easier to inspect.

I used some basic terraform to create a generic VPC with a single subnet (look at TF on how to do it or make one yourself):
This example should work fine on an intel nuc etc.  To understand how that is done take a look at the ec2-userdata.txt file in this directory.

# Step 1

With the repository cloned you can make use of the terraform code which will create a public facing VPC with a single subnet to create a single
ec2 t2.micro instance and attach the associated userdata to the instance.  Alternatively you can run the user-data script directly on an existing brand new ec2 instance


# Step 2
Validate what is happening on the ec2 instance by logging in (terraform code should output connect string)
### Note: For examples I will use Key Name: Nitecon as I use separate tf script to quickly generate, make sure you are using a Key Name that exists within your aws account.

Once logged into the system take a look at your cloudinit log details by running for example: 
```
cat /var/log/cloud-init-output.log
```


# Step 3
After terraform completed and the instance stabilizes bruce should have installed everything on the system and you should be able to hit your public innstance via a browser.

```
curl http://yourpublicIP/
```
This should give you a redirect to https://localhost, so we know it's working correctly,
Next is to view the output of the vhost that was configured with:



```
curl --header 'Host: www.example.com' --insecure https://127.0.0.1
```

Since the ssl cert we created is self signed we pass --insecure to make sure it loads.

# Step 3
Remember to clean up your test instance with something like:
```
aws ec2 terminate-instances --instance-id=i-08b84acafe719b932
```


- Enjoy
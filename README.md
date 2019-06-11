# If Tangle Then That

This small Go based application allows you to generate addresses for a 
seed and monitor those addresses for changes. If a transaction is sent to
one of those addresses it will trigger a to be defined callback. This 
allows you to easily trigger real world action based on Tangle messages.

# Settings

Settings for this application are set using environment variables;
You can either create a .env file for this or you can just set the
environment variables on your system without a .env file. 

## Available settings

 - `IFTTT_HOST`: The hostname to listen on, defaults to `localhost`
 - `IFTTT_PORT`: The port to listen on, defaults to `3693`
 - `IFTTT_SEED`: The account seed, required to run this program.
 - `IFTTT_NODE_URI`: The node to connect to, for example `https://nodes.devnet.thetangle.org:443`

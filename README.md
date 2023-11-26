Cloudflare DNS Record Creator

Simple program to automate cloudflare dns configurations.



## Getting Started
To use this tool, you first need to set up your Cloudflare API key and email. These can be provided either as environment variables (CF_API_KEY and CF_EMAIL) or in a .secrets.json file in the following format:

{
  "email": "your-email@example.com",
  "api_key": "your-api-key"
}
The tool will first try to read from the environment variables. If these are not set, it will fall back to the .secrets.json file.

## Installation

You can install the tool using the go install command:

go install github.com/<your-github-username>/<your-repo-name>@latest
Replace <your-github-username> and <your-repo-name> with your GitHub username and the name of your repository. The binary will be saved in the bin directory of your Go workspace (by default, $HOME/go/bin on Unix systems).

If your Go workspace's bin directory is in your system's PATH, you can run your application by typing its name:

<your-repo-name>

## Usage
You can run the tool with the following command-line flags:

-type: The type of the DNS record (A, CNAME, etc.)
-name: The name of the DNS record
-content: The content of the DNS record
-ttl: The TTL of the DNS record (default is 120)
-zoneid: The ID of the zone where the record will be created
The -name flag must include a TLD (e.g., dapi.nil.sd). If the domain does not include a TLD (e.g., dapi), the program will exit with an error.

If the -zoneid flag is not provided, the program will use the domain from the -name flag to retrieve the Zone ID.

Example
Here is an example of how to use the tool:

go run main.go -type=A -name=test.nil.sd -content=192.0.2.1
This command will create an A record for test.nil.sd pointing to 192.0.2.1.
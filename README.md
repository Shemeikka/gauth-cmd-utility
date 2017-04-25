# gauth-cmd-utility
Small command-line client for testing APIs.

    usage: cmd.exe [<flags>]
    
    Flags:
      --help               Show help (also see --help-long and --help-man).
      -c, --config=CONFIG  Set configuration filename without extension.

## Config-file

Example config-file with required information:

    [server]			
    	url = "http://localhost:12346"		
    	root_url = "/api/v1"		
    
    [request]
    	route = "ping"	
    	method = "POST"
    	
    	# PING
    	body = '{"ip_address": "10.12.234.139", "mac": "00:01:02:03:04:05"}'	
    
    [api]
    	[api.master]	
    	id = "public_id"
    	key = "private_key"	
    
    	[api.server]	
    	id = "public_id_serv"
    	key = "private_key_serv"

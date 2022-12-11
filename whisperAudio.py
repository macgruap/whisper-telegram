#!/usr/bin/python3

import requests, sys, whisper

r = requests.get(sys.argv[1], allow_redirects=True)
open('./audio.oga', 'wb').write(r.content)

model = whisper.load_model("large-v2")
result = model.transcribe("./audio.oga", language="Spanish")
print(result["text"])

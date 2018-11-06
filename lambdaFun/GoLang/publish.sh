\rm -rf lambda_upload.zip
zip -r lambda_upload.zip main.go
aws lambda update-function-code --function-name AudioConferenceSkill --zip-file fileb://lambda_upload.zip
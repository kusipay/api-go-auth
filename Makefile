build:
	sh ./build.sh

clean:
	rm -rf ./bin ./vendor

deploy:
	npx --yes serverless deploy

remove:
	npx --yes serverless remove

package:
	AWS_REGION=us-east-1 STAGE=qa SECOND_LEVEL_DOMAIN=kusipay TOP_LEVEL_DOMAIN=com npx --yes serverless package

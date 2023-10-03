build:
	sh ./build.sh

clean:
	rm -rf ./bin ./vendor

deploy:
	npx serverless deploy

remove:
	npx serverless remove

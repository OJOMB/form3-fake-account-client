# Fake Account API

Once the docker-compose is up the project should be testable by simply running:
```sh
make unit-test
make integration-test
```

Hopefully the integration tests will demonstrate the Modules ability to fulfill the task acceptance criteria

## A Few Things to Briefly Mention
* I wasn't exactly clear as to which Account attributes to include from those exposed by the real API. i.e. should deprecated fields be there? However as per instructions I have only left out `data.attributes.private_identification`, `data.attributes.organisation_identification` and `data.relationships`. I was hoping that worst case any extraneous fields would simply be discounted from consideration.

* I have not implemented a custom exponential back-off as recommended in the API documentation, both for the sake of expediency and because I wasn't sure if this would count as an "advanced feature", which are listed in the `Should Nots`.

* I have not included any CI, I hope this won't count against me 🤞. I am however well aware of the importance of CI and wouldn't build a production project without it
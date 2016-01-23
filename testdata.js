db.TestEnv.remove({});
db.TestEnv.insert({"AuthKey": "test-ubuntu-server-14.4x64"});
db.Task.remove({});
db.Task.insert({"AuthKey": "test-ubuntu-server-14.4x64", "PackageUrl": "github.com/regorov/logwriter"});


db.TestEnv.remove({});
db.TestEnv.insert({"authKey": "bare-metal-intel-i5"});
db.TestEnv.insert({"authKey": "digital-ocean-10$"});
db.Task.remove({});
db.Task.insert({"authKey": "bare-metal-intel-i5", "packageUrl": "github.com/regorov/logwriter"});
db.Task.insert({"authKey": "digital-ocean-10$",   "packageUrl": "github.com/regorov/logwriter"});

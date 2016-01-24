db.testEnv.remove({});
db.testEnv.insert({"authKey": "change-secret-1", "name": "bare metal (desktop)", "specification": "Intel i5, Ubuntu 14.04" });
db.testEnv.insert({"authKey": "change-secret-2", "name": "digitalocean, 10$, Frankfurt", "specification" : "Ubuntu 14.04, x64"});

db.task.remove({});
db.task.insert({"authKey": "change-secret-1", "packageName": "github.com/bradfitz/slice", "created" : new Date()});
db.task.insert({"authKey": "change-secret-1", "packageName": "github.com/valyala/fasttemplate", "created" : new Date()});
db.task.insert({"authKey": "change-secret-1", "packageName": "github.com/regorov/logwriter",    "created" : new Date()});


db.package.remove({});
db.package.insert({"name": "github.com/regorov/logwriter",
                   "url": "https://github.com/regorov/logwriter",
				   "repositoryUrl": "https://github.com",
				   "engine" : "git",
				   "created" : new Date()});

db.package.insert({"name": "github.com/nfnt/resize",
                   "url": "https://github.com/nfnt/resize",
				   "repositoryUrl": "https://github.com",
				   "engine" : "git",
				   "created" : new Date(),
				   "updated" : new Date()
				});

db.package.insert({"name": "github.com/valyala/gorpc",
                   "url": "https://github.com/valyala/gorpc",
				   "repositoryUrl": "https://github.com",
				   "engine" : "git",
				   "created" : new Date(),
				   "updated" : new Date()});

db.package.insert({"name": "github.com/valyala/fasttemplate",
                   "url": "https://github.com/valyala/fasttemplate",
				   "repositoryUrl": "https://github.com",
				   "engine" : "git",
				   "created" : new Date(),
				   "updated" : new Date()});

db.package.insert({"name": "github.com/valyala/fasthttp",
                   "url": "https://github.com/valyala/fasthttp",
				   "repositoryUrl": "https://github.com",
				   "engine" : "git",
				   "created" : new Date(),
				   "updated" : new Date()});




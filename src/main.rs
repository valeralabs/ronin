#![feature(proc_macro_hygiene, decl_macro)]

#[macro_use] extern crate rocket;
extern crate figment;
extern crate num_cpus;
#[macro_use] extern crate anyhow;

#[cfg(test)] mod tests;

mod db;
mod events;

#[get("/")]
fn hello() -> &'static str {
    "Hello, world!"
}

#[rocket::main]
async fn main() {
    // Use Rocket's default `Figment`, but allow values from `MyApp.toml`
    // and `MY_APP_` prefixed environment variables to supersede its values.
    use rocket::config::Config;

    let cpu_cores = num_cpus::get();

    if cpu_cores == 1 { panic!("You need at least 2 CPU cores to run Ronin!"); }

    // we want to share system resources between the API and the Events handler
    // events gets a single thread, and the API uses the rest
    let api_workers: usize = cpu_cores - 1;
    let event_workers: usize = 1;
    
    println!("Booting Ronin with {} API workers and {} event worker", api_workers, event_workers);
    

    // create thread for events handler
    let events_handler = std::thread::spawn( {
        // TODO: have profiles for dev, staging, prod
        // TODO: support custom profiles
        let config = Config::figment()
    		.merge(("address", "0.0.0.0"))
			.merge(("port", 3700 as usize))
			.merge(("workers", event_workers));

        rocket::custom(config)
            .mount("/", routes![hello])
            .launch()
            .await;
    });

    // main thread for API
    let config = Config::figment()
		.merge(("address", "0.0.0.0"))
		.merge(("port", 3999 as usize))
		.merge(("workers", api_workers));
        
    use db::redis::pool;

    rocket::custom(config)
        .mount("/", routes![hello])
        .manage(pool())
        .launch()
        .await;

    // wait for eventsHandler to finish before exiting
    events_handler.join().unwrap();
}


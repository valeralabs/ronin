use rocket::http;
use rocket::request;
use rocket::outcome::{Outcome, Outcome::Success, Outcome::Failure, Outcome::Forward};
use r2d2;
use r2d2_redis::RedisConnectionManager;
use anyhow::Error as AnyhowError;

const REDIS_ADDRESS: &'static str = "redis://localhost:6379";

// pool initialisation.
// called on Ronin boot for API instance + events handler
// stores a db connection pool as a rocket managed state.
pub fn pool() -> Pool {
    let manager = RedisConnectionManager::new(REDIS_ADDRESS).expect("connection manager");
	
	let pool = r2d2::Pool::builder()
		.build(manager)
		.expect("db pool");
			
	Pool(pool)
}

// Rocket guard type: a wrapper around an r2d2 pool.
// In conjunction with
//   impl<'a, 'r> request::FromRequest<'a, 'r> for RedisConnection`
// which is super confusing, it also allows code like:
//   #[post("/<item>")]
//   fn create(item: &RawStr, connection: RedisConnection) -> ...
pub struct RedisConnection(pub r2d2::PooledConnection<RedisConnectionManager>);

// alias to the type for a pool of redis connections.
pub struct Pool(pub r2d2::Pool<RedisConnectionManager>);

// retrieving a single connection from the database pool.
#[rocket::async_trait]
impl<'r> request::FromRequest<'r> for Pool {
	type Error = AnyhowError;
	
    async fn from_request(request: &'r request::Request<'_>) -> request::Outcome<Pool, Self::Error> {
        match request.guard::<Pool>().await {
        	Success(pool) => Outcome::Success(pool),
			Failure(error) => Outcome::Failure((http::Status::ServiceUnavailable, anyhow!(error.0))),
			Forward(_) => Outcome::Failure((http::Status::ServiceUnavailable, anyhow!("forward (?)"))),
        }
    }
}

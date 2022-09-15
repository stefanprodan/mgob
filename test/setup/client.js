const { MongoClient } = require("mongodb");
const PORT = process.env.MONGODB_PORT || 27017;
const HOST = process.env.MONGODB_HOST || "localhost";
const DBNAME = process.env.MONGODB_DB_NAME || "test";

async function main() {
  const client = new MongoClient(`mongodb://${HOST}:${PORT}/${DBNAME}`, {
    useUnifiedTopology: true,
  });
  try {
    await client.connect();
    await client.db().collection("test").insertOne({
      item: "item",
      val: "test",
    });
    await client.db().collection("test2").insertOne({
      item: "item2",
      val: "test2",
    });
    const result = await client.db().collection("test").findOne({});
    console.log(JSON.stringify(result));
  } catch (e) {
    console.error(e);
  } finally {
    await client.close();
  }
}

main().catch(console.error);

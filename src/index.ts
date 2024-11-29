import "dotenv/config";
import { app } from "./app";
import { PORT } from "./constants";
import db from "./database/db";

app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});

import express from "express";
import cookieParser from "cookie-parser";
import cors from "cors";
import { CROSS_ORIGIN } from "./constants";
import { ApiResponse } from "./utils/apiResponse";
import { authRouter } from "./routes/auth.routes";

const app = express();

// Allow Cross Origin Requests
app.use(
  cors({
    origin: CROSS_ORIGIN,
    credentials: true,
  })
);

// Restrict the JSON data size
app.use(express.json({ limit: "100kb" }));
// URL encoded data : extended allows to use nested objects
app.use(express.urlencoded({ extended: true, limit: "100kb" }));
// Static files
app.use(express.static("public"));
// Use Cookie Parser
app.use(cookieParser());

// Check Health
app.get("/api/v1/healthcheck", (req, res) => {
  res.status(200).json(new ApiResponse(200, "Health Check Successful", {}));
});

// Auth Routes
app.use("/api/v1/auth", authRouter);

export { app };

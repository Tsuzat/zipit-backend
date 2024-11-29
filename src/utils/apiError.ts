class ApiError extends Error {
  statusCode: number;
  message: string;
  errors: string[];

  constructor(
    statusCode: number,
    message: string,
    errors: string[] = [],
    stack: string | undefined = undefined
  ) {
    super(message);
    this.name = "ApiError";
    this.statusCode = statusCode;
    this.message = message;
    this.errors = errors;

    if (stack) {
      this.stack = stack;
    } else {
      Error.captureStackTrace(this, this.constructor);
    }
  }
}

export { ApiError };

import nodemailer from "nodemailer";
import {
  BACKEND_URL,
  EMAIL_FROM,
  EMAIL_HOST,
  EMAIL_PASSWORD,
  EMAIL_PORT,
  EMAIL_USERNAME,
} from "../constants";
import Mail from "nodemailer/lib/mailer";
import { User } from "../models/user.model";

let configOptions = {
  host: EMAIL_HOST,
  post: EMAIL_PORT,
  secure: true,
  auth: {
    user: EMAIL_USERNAME,
    pass: EMAIL_PASSWORD,
  },
};

const transporter = nodemailer.createTransport(configOptions);

const sentEmail = async (to: string, subject: string, html: string) => {
  const sendOptions: Mail.Options = {
    from: EMAIL_FROM,
    to: to,
    subject: subject,
    html: html,
  };
  transporter.sendMail(sendOptions, (error, info) => {
    if (error) {
      console.log(error);
      throw new Error(error.message);
    } else {
      console.log(`Email send to ${to} successfully. ${info.response}`);
    }
  });
};

const sendEmailVerification = async (user: User) => {
  const html = `
    <h1>Verify your email</h1>
    <p>Please click the link below to verify your email address:</p>
    <a href="${BACKEND_URL}/api/v1/auth/verify-email?email=${user.email}&verificationToken=${user.verificationToken}">
      Verify Email
    </a>
  `;
  await sentEmail(user.email, "Verify your email", html);
};

export { sentEmail, sendEmailVerification };

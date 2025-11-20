import { createAsync, query } from "@solidjs/router";
import { createTestUser, test } from "../api/api";

function wait<T>(ms: number, data: T): Promise<T> {
  return new Promise((resolve) => setTimeout(resolve, ms, data));
}

function random(min: number, max: number): number {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

const createTestUserQ = query(() => createTestUser(), "createTestUser");

const AboutData = () => {
  return createAsync(() => createTestUserQ());
};

export default AboutData;
export type AboutDataType = typeof AboutData;

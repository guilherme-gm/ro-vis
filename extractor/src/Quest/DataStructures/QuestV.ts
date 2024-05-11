import { QuestV1 } from "./QuestV1.js";
import { QuestV3 } from "./QuestV3.js";
import { QuestV4 } from "./QuestV4.js";

/**
 * Represents the different versions of quest file structures
 */
export type QuestV = QuestV1 | QuestV3 | QuestV4;

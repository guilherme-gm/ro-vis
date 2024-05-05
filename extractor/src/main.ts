import { QuestV0Parser } from "./Quest/Parsers/QuestV0Parser.js";

const parser = await QuestV0Parser.fromFile('./raw/questid2display.txt');
parser.parse();


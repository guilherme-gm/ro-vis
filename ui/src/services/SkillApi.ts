import { CommonApi, type PatchItem } from "./CommonApi";
import type { MinSkill, Skill } from "@/models/Skill";

export type SkillPatch = PatchItem<Skill>;

export class SkillApi extends CommonApi<Skill, MinSkill> {
	public static use() {
		return new SkillApi('skills/');
	}
}

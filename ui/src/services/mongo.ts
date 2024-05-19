import * as Realm from "realm-web";
import { ref } from "vue";

const app = new Realm.App({ id: "data-btpnptw" });

const appUser = ref<Realm.User | null>(null);

async function login(): Promise<Realm.User | null> {
	const user = await app.logIn(Realm.Credentials.anonymous());
	appUser.value = user;
	return user;
}

export function useMongo() {
	return {
		login,
	};
}

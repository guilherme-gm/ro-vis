import { Diff, type Change } from 'diff';

export function formatDiff(changes: Change[]) {
	let span = null;
	const fragment = document.createDocumentFragment();

	changes.forEach((part) => {
		// green for additions, red for deletions
		// grey for common parts
		const color = part.added ? 'green' :
			part.removed ? 'red' : 'grey';
		span = document.createElement('span');
		span.style.color = color;
		span.appendChild(document
			.createTextNode(part.value));
		fragment.appendChild(span);
	});

	return fragment;
}

export function diffText(from: string, to: string): DocumentFragment {
	const diff = new Diff();
	const changes = diff.diff(from, to, { ignoreCase: false });

	return formatDiff(changes);
}

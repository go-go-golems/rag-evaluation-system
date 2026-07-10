function createStore(db) {
	function nowISO() {
		return new Date().toISOString();
	}

	function slugify(title) {
		const base = String(title || "poll")
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, "-")
			.replace(/(^-+|-+$)/g, "")
			.slice(0, 40);
		return base || "poll";
	}

	function initSchema() {
		db.exec(
			"CREATE TABLE IF NOT EXISTS polls (id INTEGER PRIMARY KEY AUTOINCREMENT, slug TEXT UNIQUE, title TEXT NOT NULL, description TEXT, location TEXT, created_at TEXT NOT NULL)",
		);
		db.exec(
			"CREATE TABLE IF NOT EXISTS options (id INTEGER PRIMARY KEY AUTOINCREMENT, poll_id INTEGER NOT NULL, label TEXT NOT NULL, sort INTEGER NOT NULL)",
		);
		db.exec(
			"CREATE TABLE IF NOT EXISTS participants (id INTEGER PRIMARY KEY AUTOINCREMENT, poll_id INTEGER NOT NULL, name TEXT NOT NULL, created_at TEXT NOT NULL)",
		);
		db.exec(
			"CREATE TABLE IF NOT EXISTS votes (participant_id INTEGER NOT NULL, option_id INTEGER NOT NULL, value TEXT NOT NULL, PRIMARY KEY (participant_id, option_id))",
		);
	}

	function createPoll(title, description, location, slotLabels) {
		db.exec(
			"INSERT INTO polls (slug, title, description, location, created_at) VALUES (?, ?, ?, ?, ?)",
			null,
			title,
			description || "",
			location || "",
			nowISO(),
		);
		const pollId = Number(db.query("SELECT last_insert_rowid() AS id")[0].id);
		db.exec("UPDATE polls SET slug = ? WHERE id = ?", `${slugify(title)}-${pollId}`, pollId);
		slotLabels.forEach((label, index) => {
			db.exec("INSERT INTO options (poll_id, label, sort) VALUES (?, ?, ?)", pollId, label, index);
		});
		return pollId;
	}

	function seedIfEmpty() {
		if (Number(db.query("SELECT COUNT(*) AS n FROM polls")[0].n) !== 0) return;
		const seedId = createPoll(
			"Team offsite dinner",
			"Pick the evening that works best for the crew dinner.",
			"Trattoria Luca, downtown",
			["Thu Jul 9 · 19:00", "Fri Jul 10 · 19:00", "Sat Jul 11 · 18:30", "Sat Jul 11 · 20:00"],
		);
		const opts = db.query("SELECT id FROM options WHERE poll_id = ? ORDER BY sort", seedId);
		function seedParticipant(name, values) {
			db.exec(
				"INSERT INTO participants (poll_id, name, created_at) VALUES (?, ?, ?)",
				seedId,
				name,
				nowISO(),
			);
			const pid = Number(db.query("SELECT last_insert_rowid() AS id")[0].id);
			values.forEach((value, i) => {
				db.exec(
					"INSERT INTO votes (participant_id, option_id, value) VALUES (?, ?, ?)",
					pid,
					opts[i].id,
					value,
				);
			});
		}
		seedParticipant("Ada", ["yes", "no", "yes", "maybe"]);
		seedParticipant("Grace", ["maybe", "no", "yes", "yes"]);
		seedParticipant("Linus", ["no", "yes", "yes", "no"]);
	}

	function allPolls() {
		return db.query(
			"SELECT p.id, p.slug, p.title, p.location, p.created_at, " +
				"(SELECT COUNT(*) FROM options o WHERE o.poll_id = p.id) AS slots, " +
				"(SELECT COUNT(*) FROM participants pt WHERE pt.poll_id = p.id) AS people " +
				"FROM polls p ORDER BY p.id DESC",
		);
	}

	function getPoll(pollId) {
		const rows = db.query(
			"SELECT id, slug, title, description, location, created_at FROM polls WHERE id = ?",
			pollId,
		);
		return rows.length ? rows[0] : null;
	}

	function pollOptions(pollId) {
		return db.query(
			"SELECT id, label, sort FROM options WHERE poll_id = ? ORDER BY sort, id",
			pollId,
		);
	}

	function pollParticipants(pollId) {
		return db.query("SELECT id, name FROM participants WHERE poll_id = ? ORDER BY id", pollId);
	}

	function pollVotes(pollId) {
		return db.query(
			"SELECT v.participant_id, v.option_id, v.value FROM votes v " +
				"JOIN participants pt ON pt.id = v.participant_id WHERE pt.poll_id = ?",
			pollId,
		);
	}

	function castVote(pollId, name, body) {
		db.exec(
			"INSERT INTO participants (poll_id, name, created_at) VALUES (?, ?, ?)",
			pollId,
			name,
			nowISO(),
		);
		const participantId = Number(db.query("SELECT last_insert_rowid() AS id")[0].id);
		pollOptions(pollId).forEach((opt) => {
			const rawValue = body[`opt_${opt.id}`];
			const value =
				rawValue === "yes" || rawValue === "no" || rawValue === "maybe" ? rawValue : "no";
			db.exec(
				"INSERT INTO votes (participant_id, option_id, value) VALUES (?, ?, ?)",
				participantId,
				opt.id,
				value,
			);
		});
	}

	initSchema();
	seedIfEmpty();

	return {
		nowISO,
		createPoll,
		allPolls,
		getPoll,
		pollOptions,
		pollParticipants,
		pollVotes,
		castVote,
	};
}

module.exports = { createStore };

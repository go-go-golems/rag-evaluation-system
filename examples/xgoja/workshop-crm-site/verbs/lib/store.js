// SQLite persistence for the workshop CRM reference host. CRM definitions live
// in pages.js; records below are plain, serializable application data.
function createStore(db) {
	function nowISO() {
		return new Date().toISOString();
	}

	function initSchema() {
		db.exec(
			"CREATE TABLE IF NOT EXISTS organizations (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, industry TEXT, created_at TEXT NOT NULL)",
		);
		db.exec(
			"CREATE TABLE IF NOT EXISTS contacts (id INTEGER PRIMARY KEY AUTOINCREMENT, organization_id INTEGER NOT NULL, name TEXT NOT NULL, email TEXT NOT NULL, created_at TEXT NOT NULL)",
		);
		db.exec(
			"CREATE TABLE IF NOT EXISTS deals (id INTEGER PRIMARY KEY AUTOINCREMENT, organization_id INTEGER NOT NULL, contact_id INTEGER NOT NULL, title TEXT NOT NULL, stage_id TEXT NOT NULL, amount INTEGER NOT NULL, format TEXT NOT NULL, workshop_date TEXT, created_at TEXT NOT NULL)",
		);
		db.exec(
			"CREATE TABLE IF NOT EXISTS activities (id INTEGER PRIMARY KEY AUTOINCREMENT, deal_id INTEGER NOT NULL, kind TEXT NOT NULL, title TEXT NOT NULL, at_iso TEXT NOT NULL)",
		);
		db.exec(
			"CREATE TABLE IF NOT EXISTS availability_options (id INTEGER PRIMARY KEY AUTOINCREMENT, deal_id INTEGER NOT NULL, label TEXT NOT NULL, start_iso TEXT NOT NULL, end_iso TEXT NOT NULL, sort INTEGER NOT NULL)",
		);
		db.exec(
			"CREATE TABLE IF NOT EXISTS availability_votes (option_id INTEGER NOT NULL, name TEXT NOT NULL, value TEXT NOT NULL, PRIMARY KEY (option_id, name))",
		);
		db.exec(
			"CREATE TABLE IF NOT EXISTS workshop_runs (id INTEGER PRIMARY KEY AUTOINCREMENT, deal_id INTEGER NOT NULL UNIQUE, title TEXT NOT NULL, start_iso TEXT NOT NULL, end_iso TEXT NOT NULL, status TEXT NOT NULL, created_at TEXT NOT NULL)",
		);
	}

	function idAfterInsert() {
		return Number(db.query("SELECT last_insert_rowid() AS id")[0].id);
	}

	function addActivity(dealId, kind, title) {
		db.exec(
			"INSERT INTO activities (deal_id, kind, title, at_iso) VALUES (?, ?, ?, ?)",
			dealId,
			kind,
			title,
			nowISO(),
		);
	}

	function createLead({ organization, contact, email, amount, format }) {
		const createdAt = nowISO();
		db.exec(
			"INSERT INTO organizations (name, industry, created_at) VALUES (?, ?, ?)",
			organization,
			"Technology",
			createdAt,
		);
		const organizationId = idAfterInsert();
		db.exec(
			"INSERT INTO contacts (organization_id, name, email, created_at) VALUES (?, ?, ?, ?)",
			organizationId,
			contact,
			email,
			createdAt,
		);
		const contactId = idAfterInsert();
		db.exec(
			"INSERT INTO deals (organization_id, contact_id, title, stage_id, amount, format, workshop_date, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			organizationId,
			contactId,
			`${organization} AI engineering workshop`,
			"lead",
			Number(amount || 0),
			format || "onsite-2d",
			null,
			createdAt,
		);
		const dealId = idAfterInsert();
		addActivity(dealId, "note", "Lead captured from the workshop intake form");
		return dealId;
	}

	function seedIfEmpty() {
		if (Number(db.query("SELECT COUNT(*) AS n FROM deals")[0].n) > 0) return;
		const dealId = createLead({
			organization: "Acme Robotics",
			contact: "Maya Chen",
			email: "maya@acme.example",
			amount: 18000,
			format: "onsite-2d",
		});
		db.exec("UPDATE deals SET stage_id = ? WHERE id = ?", "proposal", dealId);
		addActivity(dealId, "email", "Sent the draft two-day AI engineering agenda");
		addActivity(
			dealId,
			"meeting",
			"Discovery call: platform team wants retrieval evaluation exercises",
		);
		[
			["2026-08-18 · 09:00", "2026-08-18T09:00:00", "2026-08-18T17:00:00"],
			["2026-08-20 · 09:00", "2026-08-20T09:00:00", "2026-08-20T17:00:00"],
		].forEach((slot, index) => {
			db.exec(
				"INSERT INTO availability_options (deal_id, label, start_iso, end_iso, sort) VALUES (?, ?, ?, ?, ?)",
				dealId,
				slot[0],
				slot[1],
				slot[2],
				index,
			);
		});
	}

	function allDeals() {
		return db.query(
			"SELECT d.id, d.title, d.stage_id AS stageId, d.amount, d.format, d.workshop_date AS workshopDate, o.name AS organization, c.name AS contact, c.email FROM deals d JOIN organizations o ON o.id = d.organization_id JOIN contacts c ON c.id = d.contact_id ORDER BY d.id DESC",
		);
	}

	function getDeal(dealId) {
		const rows = db.query(
			"SELECT d.id, d.title, d.stage_id AS stageId, d.amount, d.format, d.workshop_date AS workshopDate, o.name AS organization, c.name AS contact, c.email FROM deals d JOIN organizations o ON o.id = d.organization_id JOIN contacts c ON c.id = d.contact_id WHERE d.id = ?",
			dealId,
		);
		return rows.length ? rows[0] : null;
	}

	function activitiesForDeal(dealId) {
		return db.query(
			"SELECT id, kind, title, at_iso AS atISO FROM activities WHERE deal_id = ? ORDER BY id DESC",
			dealId,
		);
	}

	function availabilityForDeal(dealId) {
		return db.query(
			"SELECT id, label, start_iso AS startISO, end_iso AS endISO FROM availability_options WHERE deal_id = ? ORDER BY sort, id",
			dealId,
		);
	}

	function moveDeal(dealId, toStage) {
		const validStages = ["lead", "proposal", "won"];
		if (!validStages.includes(toStage) || !getDeal(dealId)) return false;
		db.exec("UPDATE deals SET stage_id = ? WHERE id = ?", toStage, dealId);
		addActivity(dealId, "stage_change", `Moved opportunity to ${toStage}`);
		return true;
	}

	function createAvailability(dealId) {
		if (availabilityForDeal(dealId).length) return;
		[
			["2026-09-08 · 09:00", "2026-09-08T09:00:00", "2026-09-08T17:00:00"],
			["2026-09-10 · 09:00", "2026-09-10T09:00:00", "2026-09-10T17:00:00"],
		].forEach((slot, index) => {
			db.exec(
				"INSERT INTO availability_options (deal_id, label, start_iso, end_iso, sort) VALUES (?, ?, ?, ?, ?)",
				dealId,
				slot[0],
				slot[1],
				slot[2],
				index,
			);
		});
		addActivity(dealId, "note", "Created a delivery-date availability poll");
	}

	function scheduleRun(dealId, optionId) {
		const options = db.query(
			"SELECT id, label, start_iso AS startISO, end_iso AS endISO FROM availability_options WHERE id = ? AND deal_id = ?",
			optionId,
			dealId,
		);
		if (!options.length) return false;
		const option = options[0];
		const deal = getDeal(dealId);
		db.exec(
			"INSERT OR REPLACE INTO workshop_runs (deal_id, title, start_iso, end_iso, status, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			dealId,
			deal.title,
			option.startISO,
			option.endISO,
			"scheduled",
			nowISO(),
		);
		db.exec(
			"UPDATE deals SET stage_id = ?, workshop_date = ? WHERE id = ?",
			"won",
			option.startISO.slice(0, 10),
			dealId,
		);
		addActivity(dealId, "meeting", `Scheduled workshop for ${option.label}`);
		return true;
	}

	function allRuns() {
		return db.query(
			"SELECT id, deal_id AS dealId, title, start_iso AS startISO, end_iso AS endISO, status FROM workshop_runs ORDER BY start_iso",
		);
	}

	initSchema();
	seedIfEmpty();
	return {
		allDeals,
		getDeal,
		activitiesForDeal,
		availabilityForDeal,
		createLead,
		moveDeal,
		createAvailability,
		scheduleRun,
		allRuns,
	};
}

module.exports = { createStore };

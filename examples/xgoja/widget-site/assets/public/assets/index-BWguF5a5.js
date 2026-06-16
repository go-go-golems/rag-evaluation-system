var e = (e, t) => () => (t || (e((t = { exports: {} }).exports, t), (e = null)), t.exports);
(function () {
	let e = document.createElement(`link`).relList;
	if (e && e.supports && e.supports(`modulepreload`)) return;
	for (let e of document.querySelectorAll(`link[rel="modulepreload"]`)) n(e);
	new MutationObserver((e) => {
		for (let t of e)
			if (t.type === `childList`)
				for (let e of t.addedNodes) e.tagName === `LINK` && e.rel === `modulepreload` && n(e);
	}).observe(document, { childList: !0, subtree: !0 });
	function t(e) {
		let t = {};
		return (
			e.integrity && (t.integrity = e.integrity),
			e.referrerPolicy && (t.referrerPolicy = e.referrerPolicy),
			e.crossOrigin === `use-credentials`
				? (t.credentials = `include`)
				: e.crossOrigin === `anonymous`
					? (t.credentials = `omit`)
					: (t.credentials = `same-origin`),
			t
		);
	}
	function n(e) {
		if (e.ep) return;
		e.ep = !0;
		let n = t(e);
		fetch(e.href, n);
	}
})();
var t = e((e) => {
		var t = Symbol.for(`react.transitional.element`),
			n = Symbol.for(`react.portal`),
			r = Symbol.for(`react.fragment`),
			i = Symbol.for(`react.strict_mode`),
			a = Symbol.for(`react.profiler`),
			o = Symbol.for(`react.consumer`),
			s = Symbol.for(`react.context`),
			c = Symbol.for(`react.forward_ref`),
			l = Symbol.for(`react.suspense`),
			u = Symbol.for(`react.memo`),
			d = Symbol.for(`react.lazy`),
			f = Symbol.for(`react.activity`),
			p = Symbol.iterator;
		function m(e) {
			return typeof e != `object` || !e
				? null
				: ((e = (p && e[p]) || e[`@@iterator`]), typeof e == `function` ? e : null);
		}
		var h = {
				isMounted: function () {
					return !1;
				},
				enqueueForceUpdate: function () {},
				enqueueReplaceState: function () {},
				enqueueSetState: function () {},
			},
			g = Object.assign,
			_ = {};
		function v(e, t, n) {
			(this.props = e), (this.context = t), (this.refs = _), (this.updater = n || h);
		}
		(v.prototype.isReactComponent = {}),
			(v.prototype.setState = function (e, t) {
				if (typeof e != `object` && typeof e != `function` && e != null)
					throw Error(
						`takes an object of state variables to update or a function which returns an object of state variables.`,
					);
				this.updater.enqueueSetState(this, e, t, `setState`);
			}),
			(v.prototype.forceUpdate = function (e) {
				this.updater.enqueueForceUpdate(this, e, `forceUpdate`);
			});
		function y() {}
		y.prototype = v.prototype;
		function b(e, t, n) {
			(this.props = e), (this.context = t), (this.refs = _), (this.updater = n || h);
		}
		var x = (b.prototype = new y());
		(x.constructor = b), g(x, v.prototype), (x.isPureReactComponent = !0);
		var ee = Array.isArray;
		function S() {}
		var C = { H: null, A: null, T: null, S: null },
			te = Object.prototype.hasOwnProperty;
		function ne(e, n, r) {
			var i = r.ref;
			return { $$typeof: t, type: e, key: n, ref: i === void 0 ? null : i, props: r };
		}
		function re(e, t) {
			return ne(e.type, t, e.props);
		}
		function w(e) {
			return typeof e == `object` && !!e && e.$$typeof === t;
		}
		function ie(e) {
			var t = { "=": `=0`, ":": `=2` };
			return (
				`$` +
				e.replace(/[=:]/g, function (e) {
					return t[e];
				})
			);
		}
		var ae = /\/+/g;
		function oe(e, t) {
			return typeof e == `object` && e && e.key != null ? ie(`` + e.key) : t.toString(36);
		}
		function se(e) {
			switch (e.status) {
				case `fulfilled`:
					return e.value;
				case `rejected`:
					throw e.reason;
				default:
					switch (
						(typeof e.status == `string`
							? e.then(S, S)
							: ((e.status = `pending`),
								e.then(
									function (t) {
										e.status === `pending` && ((e.status = `fulfilled`), (e.value = t));
									},
									function (t) {
										e.status === `pending` && ((e.status = `rejected`), (e.reason = t));
									},
								)),
						e.status)
					) {
						case `fulfilled`:
							return e.value;
						case `rejected`:
							throw e.reason;
					}
			}
			throw e;
		}
		function ce(e, r, i, a, o) {
			var s = typeof e;
			(s === `undefined` || s === `boolean`) && (e = null);
			var c = !1;
			if (e === null) c = !0;
			else
				switch (s) {
					case `bigint`:
					case `string`:
					case `number`:
						c = !0;
						break;
					case `object`:
						switch (e.$$typeof) {
							case t:
							case n:
								c = !0;
								break;
							case d:
								return (c = e._init), ce(c(e._payload), r, i, a, o);
						}
				}
			if (c)
				return (
					(o = o(e)),
					(c = a === `` ? `.` + oe(e, 0) : a),
					ee(o)
						? ((i = ``),
							c != null && (i = c.replace(ae, `$&/`) + `/`),
							ce(o, r, i, ``, function (e) {
								return e;
							}))
						: o != null &&
							(w(o) &&
								(o = re(
									o,
									i +
										(o.key == null || (e && e.key === o.key)
											? ``
											: (`` + o.key).replace(ae, `$&/`) + `/`) +
										c,
								)),
							r.push(o)),
					1
				);
			c = 0;
			var l = a === `` ? `.` : a + `:`;
			if (ee(e))
				for (var u = 0; u < e.length; u++) (a = e[u]), (s = l + oe(a, u)), (c += ce(a, r, i, s, o));
			else if (((u = m(e)), typeof u == `function`))
				for (e = u.call(e), u = 0; !(a = e.next()).done; )
					(a = a.value), (s = l + oe(a, u++)), (c += ce(a, r, i, s, o));
			else if (s === `object`) {
				if (typeof e.then == `function`) return ce(se(e), r, i, a, o);
				throw (
					((r = String(e)),
					Error(
						`Objects are not valid as a React child (found: ` +
							(r === `[object Object]`
								? `object with keys {` + Object.keys(e).join(`, `) + `}`
								: r) +
							`). If you meant to render a collection of children, use an array instead.`,
					))
				);
			}
			return c;
		}
		function le(e, t, n) {
			if (e == null) return e;
			var r = [],
				i = 0;
			return (
				ce(e, r, ``, ``, function (e) {
					return t.call(n, e, i++);
				}),
				r
			);
		}
		function ue(e) {
			if (e._status === -1) {
				var t = e._result;
				(t = t()),
					t.then(
						function (t) {
							(e._status === 0 || e._status === -1) && ((e._status = 1), (e._result = t));
						},
						function (t) {
							(e._status === 0 || e._status === -1) && ((e._status = 2), (e._result = t));
						},
					),
					e._status === -1 && ((e._status = 0), (e._result = t));
			}
			if (e._status === 1) return e._result.default;
			throw e._result;
		}
		var T =
				typeof reportError == `function`
					? reportError
					: function (e) {
							if (typeof window == `object` && typeof window.ErrorEvent == `function`) {
								var t = new window.ErrorEvent(`error`, {
									bubbles: !0,
									cancelable: !0,
									message:
										typeof e == `object` && e && typeof e.message == `string`
											? String(e.message)
											: String(e),
									error: e,
								});
								if (!window.dispatchEvent(t)) return;
							} else if (typeof process == `object` && typeof process.emit == `function`) {
								process.emit(`uncaughtException`, e);
								return;
							}
							console.error(e);
						},
			E = {
				map: le,
				forEach: function (e, t, n) {
					le(
						e,
						function () {
							t.apply(this, arguments);
						},
						n,
					);
				},
				count: function (e) {
					var t = 0;
					return (
						le(e, function () {
							t++;
						}),
						t
					);
				},
				toArray: function (e) {
					return (
						le(e, function (e) {
							return e;
						}) || []
					);
				},
				only: function (e) {
					if (!w(e))
						throw Error(`React.Children.only expected to receive a single React element child.`);
					return e;
				},
			};
		(e.Activity = f),
			(e.Children = E),
			(e.Component = v),
			(e.Fragment = r),
			(e.Profiler = a),
			(e.PureComponent = b),
			(e.StrictMode = i),
			(e.Suspense = l),
			(e.__CLIENT_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE = C),
			(e.__COMPILER_RUNTIME = {
				__proto__: null,
				c: function (e) {
					return C.H.useMemoCache(e);
				},
			}),
			(e.cache = function (e) {
				return function () {
					return e.apply(null, arguments);
				};
			}),
			(e.cacheSignal = function () {
				return null;
			}),
			(e.cloneElement = function (e, t, n) {
				if (e == null)
					throw Error(`The argument must be a React element, but you passed ` + e + `.`);
				var r = g({}, e.props),
					i = e.key;
				if (t != null)
					for (a in (t.key !== void 0 && (i = `` + t.key), t))
						!te.call(t, a) ||
							a === `key` ||
							a === `__self` ||
							a === `__source` ||
							(a === `ref` && t.ref === void 0) ||
							(r[a] = t[a]);
				var a = arguments.length - 2;
				if (a === 1) r.children = n;
				else if (1 < a) {
					for (var o = Array(a), s = 0; s < a; s++) o[s] = arguments[s + 2];
					r.children = o;
				}
				return ne(e.type, i, r);
			}),
			(e.createContext = function (e) {
				return (
					(e = {
						$$typeof: s,
						_currentValue: e,
						_currentValue2: e,
						_threadCount: 0,
						Provider: null,
						Consumer: null,
					}),
					(e.Provider = e),
					(e.Consumer = { $$typeof: o, _context: e }),
					e
				);
			}),
			(e.createElement = function (e, t, n) {
				var r,
					i = {},
					a = null;
				if (t != null)
					for (r in (t.key !== void 0 && (a = `` + t.key), t))
						te.call(t, r) && r !== `key` && r !== `__self` && r !== `__source` && (i[r] = t[r]);
				var o = arguments.length - 2;
				if (o === 1) i.children = n;
				else if (1 < o) {
					for (var s = Array(o), c = 0; c < o; c++) s[c] = arguments[c + 2];
					i.children = s;
				}
				if (e && e.defaultProps)
					for (r in ((o = e.defaultProps), o)) i[r] === void 0 && (i[r] = o[r]);
				return ne(e, a, i);
			}),
			(e.createRef = function () {
				return { current: null };
			}),
			(e.forwardRef = function (e) {
				return { $$typeof: c, render: e };
			}),
			(e.isValidElement = w),
			(e.lazy = function (e) {
				return { $$typeof: d, _payload: { _status: -1, _result: e }, _init: ue };
			}),
			(e.memo = function (e, t) {
				return { $$typeof: u, type: e, compare: t === void 0 ? null : t };
			}),
			(e.startTransition = function (e) {
				var t = C.T,
					n = {};
				C.T = n;
				try {
					var r = e(),
						i = C.S;
					i !== null && i(n, r),
						typeof r == `object` && r && typeof r.then == `function` && r.then(S, T);
				} catch (e) {
					T(e);
				} finally {
					t !== null && n.types !== null && (t.types = n.types), (C.T = t);
				}
			}),
			(e.unstable_useCacheRefresh = function () {
				return C.H.useCacheRefresh();
			}),
			(e.use = function (e) {
				return C.H.use(e);
			}),
			(e.useActionState = function (e, t, n) {
				return C.H.useActionState(e, t, n);
			}),
			(e.useCallback = function (e, t) {
				return C.H.useCallback(e, t);
			}),
			(e.useContext = function (e) {
				return C.H.useContext(e);
			}),
			(e.useDebugValue = function () {}),
			(e.useDeferredValue = function (e, t) {
				return C.H.useDeferredValue(e, t);
			}),
			(e.useEffect = function (e, t) {
				return C.H.useEffect(e, t);
			}),
			(e.useEffectEvent = function (e) {
				return C.H.useEffectEvent(e);
			}),
			(e.useId = function () {
				return C.H.useId();
			}),
			(e.useImperativeHandle = function (e, t, n) {
				return C.H.useImperativeHandle(e, t, n);
			}),
			(e.useInsertionEffect = function (e, t) {
				return C.H.useInsertionEffect(e, t);
			}),
			(e.useLayoutEffect = function (e, t) {
				return C.H.useLayoutEffect(e, t);
			}),
			(e.useMemo = function (e, t) {
				return C.H.useMemo(e, t);
			}),
			(e.useOptimistic = function (e, t) {
				return C.H.useOptimistic(e, t);
			}),
			(e.useReducer = function (e, t, n) {
				return C.H.useReducer(e, t, n);
			}),
			(e.useRef = function (e) {
				return C.H.useRef(e);
			}),
			(e.useState = function (e) {
				return C.H.useState(e);
			}),
			(e.useSyncExternalStore = function (e, t, n) {
				return C.H.useSyncExternalStore(e, t, n);
			}),
			(e.useTransition = function () {
				return C.H.useTransition();
			}),
			(e.version = `19.2.7`);
	}),
	n = e((e, n) => {
		n.exports = t();
	}),
	r = e((e) => {
		function t(e, t) {
			var n = e.length;
			e.push(t);
			a: for (; 0 < n; ) {
				var r = (n - 1) >>> 1,
					a = e[r];
				if (0 < i(a, t)) (e[r] = t), (e[n] = a), (n = r);
				else break a;
			}
		}
		function n(e) {
			return e.length === 0 ? null : e[0];
		}
		function r(e) {
			if (e.length === 0) return null;
			var t = e[0],
				n = e.pop();
			if (n !== t) {
				e[0] = n;
				a: for (var r = 0, a = e.length, o = a >>> 1; r < o; ) {
					var s = 2 * (r + 1) - 1,
						c = e[s],
						l = s + 1,
						u = e[l];
					if (0 > i(c, n))
						l < a && 0 > i(u, c)
							? ((e[r] = u), (e[l] = n), (r = l))
							: ((e[r] = c), (e[s] = n), (r = s));
					else if (l < a && 0 > i(u, n)) (e[r] = u), (e[l] = n), (r = l);
					else break a;
				}
			}
			return t;
		}
		function i(e, t) {
			var n = e.sortIndex - t.sortIndex;
			return n === 0 ? e.id - t.id : n;
		}
		if (
			((e.unstable_now = void 0),
			typeof performance == `object` && typeof performance.now == `function`)
		) {
			var a = performance;
			e.unstable_now = function () {
				return a.now();
			};
		} else {
			var o = Date,
				s = o.now();
			e.unstable_now = function () {
				return o.now() - s;
			};
		}
		var c = [],
			l = [],
			u = 1,
			d = null,
			f = 3,
			p = !1,
			m = !1,
			h = !1,
			g = !1,
			_ = typeof setTimeout == `function` ? setTimeout : null,
			v = typeof clearTimeout == `function` ? clearTimeout : null,
			y = typeof setImmediate < `u` ? setImmediate : null;
		function b(e) {
			for (var i = n(l); i !== null; ) {
				if (i.callback === null) r(l);
				else if (i.startTime <= e) r(l), (i.sortIndex = i.expirationTime), t(c, i);
				else break;
				i = n(l);
			}
		}
		function x(e) {
			if (((h = !1), b(e), !m))
				if (n(c) !== null) (m = !0), ee || ((ee = !0), w());
				else {
					var t = n(l);
					t !== null && oe(x, t.startTime - e);
				}
		}
		var ee = !1,
			S = -1,
			C = 5,
			te = -1;
		function ne() {
			return g ? !0 : !(e.unstable_now() - te < C);
		}
		function re() {
			if (((g = !1), ee)) {
				var t = e.unstable_now();
				te = t;
				var i = !0;
				try {
					a: {
						(m = !1), h && ((h = !1), v(S), (S = -1)), (p = !0);
						var a = f;
						try {
							b: {
								for (b(t), d = n(c); d !== null && !(d.expirationTime > t && ne()); ) {
									var o = d.callback;
									if (typeof o == `function`) {
										(d.callback = null), (f = d.priorityLevel);
										var s = o(d.expirationTime <= t);
										if (((t = e.unstable_now()), typeof s == `function`)) {
											(d.callback = s), b(t), (i = !0);
											break b;
										}
										d === n(c) && r(c), b(t);
									} else r(c);
									d = n(c);
								}
								if (d !== null) i = !0;
								else {
									var u = n(l);
									u !== null && oe(x, u.startTime - t), (i = !1);
								}
							}
							break a;
						} finally {
							(d = null), (f = a), (p = !1);
						}
						i = void 0;
					}
				} finally {
					i ? w() : (ee = !1);
				}
			}
		}
		var w;
		if (typeof y == `function`)
			w = function () {
				y(re);
			};
		else if (typeof MessageChannel < `u`) {
			var ie = new MessageChannel(),
				ae = ie.port2;
			(ie.port1.onmessage = re),
				(w = function () {
					ae.postMessage(null);
				});
		} else
			w = function () {
				_(re, 0);
			};
		function oe(t, n) {
			S = _(function () {
				t(e.unstable_now());
			}, n);
		}
		(e.unstable_IdlePriority = 5),
			(e.unstable_ImmediatePriority = 1),
			(e.unstable_LowPriority = 4),
			(e.unstable_NormalPriority = 3),
			(e.unstable_Profiling = null),
			(e.unstable_UserBlockingPriority = 2),
			(e.unstable_cancelCallback = function (e) {
				e.callback = null;
			}),
			(e.unstable_forceFrameRate = function (e) {
				0 > e || 125 < e
					? console.error(
							`forceFrameRate takes a positive int between 0 and 125, forcing frame rates higher than 125 fps is not supported`,
						)
					: (C = 0 < e ? Math.floor(1e3 / e) : 5);
			}),
			(e.unstable_getCurrentPriorityLevel = function () {
				return f;
			}),
			(e.unstable_next = function (e) {
				switch (f) {
					case 1:
					case 2:
					case 3:
						var t = 3;
						break;
					default:
						t = f;
				}
				var n = f;
				f = t;
				try {
					return e();
				} finally {
					f = n;
				}
			}),
			(e.unstable_requestPaint = function () {
				g = !0;
			}),
			(e.unstable_runWithPriority = function (e, t) {
				switch (e) {
					case 1:
					case 2:
					case 3:
					case 4:
					case 5:
						break;
					default:
						e = 3;
				}
				var n = f;
				f = e;
				try {
					return t();
				} finally {
					f = n;
				}
			}),
			(e.unstable_scheduleCallback = function (r, i, a) {
				var o = e.unstable_now();
				switch (
					(typeof a == `object` && a
						? ((a = a.delay), (a = typeof a == `number` && 0 < a ? o + a : o))
						: (a = o),
					r)
				) {
					case 1:
						var s = -1;
						break;
					case 2:
						s = 250;
						break;
					case 5:
						s = 1073741823;
						break;
					case 4:
						s = 1e4;
						break;
					default:
						s = 5e3;
				}
				return (
					(s = a + s),
					(r = {
						id: u++,
						callback: i,
						priorityLevel: r,
						startTime: a,
						expirationTime: s,
						sortIndex: -1,
					}),
					a > o
						? ((r.sortIndex = a),
							t(l, r),
							n(c) === null && r === n(l) && (h ? (v(S), (S = -1)) : (h = !0), oe(x, a - o)))
						: ((r.sortIndex = s), t(c, r), m || p || ((m = !0), ee || ((ee = !0), w()))),
					r
				);
			}),
			(e.unstable_shouldYield = ne),
			(e.unstable_wrapCallback = function (e) {
				var t = f;
				return function () {
					var n = f;
					f = t;
					try {
						return e.apply(this, arguments);
					} finally {
						f = n;
					}
				};
			});
	}),
	i = e((e, t) => {
		t.exports = r();
	}),
	a = e((e) => {
		var t = n();
		function r(e) {
			var t = `https://react.dev/errors/` + e;
			if (1 < arguments.length) {
				t += `?args[]=` + encodeURIComponent(arguments[1]);
				for (var n = 2; n < arguments.length; n++)
					t += `&args[]=` + encodeURIComponent(arguments[n]);
			}
			return (
				`Minified React error #` +
				e +
				`; visit ` +
				t +
				` for the full message or use the non-minified dev environment for full errors and additional helpful warnings.`
			);
		}
		function i() {}
		var a = {
				d: {
					f: i,
					r: function () {
						throw Error(r(522));
					},
					D: i,
					C: i,
					L: i,
					m: i,
					X: i,
					S: i,
					M: i,
				},
				p: 0,
				findDOMNode: null,
			},
			o = Symbol.for(`react.portal`);
		function s(e, t, n) {
			var r = 3 < arguments.length && arguments[3] !== void 0 ? arguments[3] : null;
			return {
				$$typeof: o,
				key: r == null ? null : `` + r,
				children: e,
				containerInfo: t,
				implementation: n,
			};
		}
		var c = t.__CLIENT_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE;
		function l(e, t) {
			if (e === `font`) return ``;
			if (typeof t == `string`) return t === `use-credentials` ? t : ``;
		}
		(e.__DOM_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE = a),
			(e.createPortal = function (e, t) {
				var n = 2 < arguments.length && arguments[2] !== void 0 ? arguments[2] : null;
				if (!t || (t.nodeType !== 1 && t.nodeType !== 9 && t.nodeType !== 11)) throw Error(r(299));
				return s(e, t, null, n);
			}),
			(e.flushSync = function (e) {
				var t = c.T,
					n = a.p;
				try {
					if (((c.T = null), (a.p = 2), e)) return e();
				} finally {
					(c.T = t), (a.p = n), a.d.f();
				}
			}),
			(e.preconnect = function (e, t) {
				typeof e == `string` &&
					(t
						? ((t = t.crossOrigin),
							(t = typeof t == `string` ? (t === `use-credentials` ? t : ``) : void 0))
						: (t = null),
					a.d.C(e, t));
			}),
			(e.prefetchDNS = function (e) {
				typeof e == `string` && a.d.D(e);
			}),
			(e.preinit = function (e, t) {
				if (typeof e == `string` && t && typeof t.as == `string`) {
					var n = t.as,
						r = l(n, t.crossOrigin),
						i = typeof t.integrity == `string` ? t.integrity : void 0,
						o = typeof t.fetchPriority == `string` ? t.fetchPriority : void 0;
					n === `style`
						? a.d.S(e, typeof t.precedence == `string` ? t.precedence : void 0, {
								crossOrigin: r,
								integrity: i,
								fetchPriority: o,
							})
						: n === `script` &&
							a.d.X(e, {
								crossOrigin: r,
								integrity: i,
								fetchPriority: o,
								nonce: typeof t.nonce == `string` ? t.nonce : void 0,
							});
				}
			}),
			(e.preinitModule = function (e, t) {
				if (typeof e == `string`)
					if (typeof t == `object` && t) {
						if (t.as == null || t.as === `script`) {
							var n = l(t.as, t.crossOrigin);
							a.d.M(e, {
								crossOrigin: n,
								integrity: typeof t.integrity == `string` ? t.integrity : void 0,
								nonce: typeof t.nonce == `string` ? t.nonce : void 0,
							});
						}
					} else t ?? a.d.M(e);
			}),
			(e.preload = function (e, t) {
				if (typeof e == `string` && typeof t == `object` && t && typeof t.as == `string`) {
					var n = t.as,
						r = l(n, t.crossOrigin);
					a.d.L(e, n, {
						crossOrigin: r,
						integrity: typeof t.integrity == `string` ? t.integrity : void 0,
						nonce: typeof t.nonce == `string` ? t.nonce : void 0,
						type: typeof t.type == `string` ? t.type : void 0,
						fetchPriority: typeof t.fetchPriority == `string` ? t.fetchPriority : void 0,
						referrerPolicy: typeof t.referrerPolicy == `string` ? t.referrerPolicy : void 0,
						imageSrcSet: typeof t.imageSrcSet == `string` ? t.imageSrcSet : void 0,
						imageSizes: typeof t.imageSizes == `string` ? t.imageSizes : void 0,
						media: typeof t.media == `string` ? t.media : void 0,
					});
				}
			}),
			(e.preloadModule = function (e, t) {
				if (typeof e == `string`)
					if (t) {
						var n = l(t.as, t.crossOrigin);
						a.d.m(e, {
							as: typeof t.as == `string` && t.as !== `script` ? t.as : void 0,
							crossOrigin: n,
							integrity: typeof t.integrity == `string` ? t.integrity : void 0,
						});
					} else a.d.m(e);
			}),
			(e.requestFormReset = function (e) {
				a.d.r(e);
			}),
			(e.unstable_batchedUpdates = function (e, t) {
				return e(t);
			}),
			(e.useFormState = function (e, t, n) {
				return c.H.useFormState(e, t, n);
			}),
			(e.useFormStatus = function () {
				return c.H.useHostTransitionStatus();
			}),
			(e.version = `19.2.7`);
	}),
	o = e((e, t) => {
		function n() {
			if (
				!(
					typeof __REACT_DEVTOOLS_GLOBAL_HOOK__ > `u` ||
					typeof __REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE != `function`
				)
			)
				try {
					__REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE(n);
				} catch (e) {
					console.error(e);
				}
		}
		n(), (t.exports = a());
	}),
	s = e((e) => {
		var t = i(),
			r = n(),
			a = o();
		function s(e) {
			var t = `https://react.dev/errors/` + e;
			if (1 < arguments.length) {
				t += `?args[]=` + encodeURIComponent(arguments[1]);
				for (var n = 2; n < arguments.length; n++)
					t += `&args[]=` + encodeURIComponent(arguments[n]);
			}
			return (
				`Minified React error #` +
				e +
				`; visit ` +
				t +
				` for the full message or use the non-minified dev environment for full errors and additional helpful warnings.`
			);
		}
		function c(e) {
			return !(!e || (e.nodeType !== 1 && e.nodeType !== 9 && e.nodeType !== 11));
		}
		function l(e) {
			var t = e,
				n = e;
			if (e.alternate) for (; t.return; ) t = t.return;
			else {
				e = t;
				do (t = e), t.flags & 4098 && (n = t.return), (e = t.return);
				while (e);
			}
			return t.tag === 3 ? n : null;
		}
		function u(e) {
			if (e.tag === 13) {
				var t = e.memoizedState;
				if ((t === null && ((e = e.alternate), e !== null && (t = e.memoizedState)), t !== null))
					return t.dehydrated;
			}
			return null;
		}
		function d(e) {
			if (e.tag === 31) {
				var t = e.memoizedState;
				if ((t === null && ((e = e.alternate), e !== null && (t = e.memoizedState)), t !== null))
					return t.dehydrated;
			}
			return null;
		}
		function f(e) {
			if (l(e) !== e) throw Error(s(188));
		}
		function p(e) {
			var t = e.alternate;
			if (!t) {
				if (((t = l(e)), t === null)) throw Error(s(188));
				return t === e ? e : null;
			}
			for (var n = e, r = t; ; ) {
				var i = n.return;
				if (i === null) break;
				var a = i.alternate;
				if (a === null) {
					if (((r = i.return), r !== null)) {
						n = r;
						continue;
					}
					break;
				}
				if (i.child === a.child) {
					for (a = i.child; a; ) {
						if (a === n) return f(i), e;
						if (a === r) return f(i), t;
						a = a.sibling;
					}
					throw Error(s(188));
				}
				if (n.return !== r.return) (n = i), (r = a);
				else {
					for (var o = !1, c = i.child; c; ) {
						if (c === n) {
							(o = !0), (n = i), (r = a);
							break;
						}
						if (c === r) {
							(o = !0), (r = i), (n = a);
							break;
						}
						c = c.sibling;
					}
					if (!o) {
						for (c = a.child; c; ) {
							if (c === n) {
								(o = !0), (n = a), (r = i);
								break;
							}
							if (c === r) {
								(o = !0), (r = a), (n = i);
								break;
							}
							c = c.sibling;
						}
						if (!o) throw Error(s(189));
					}
				}
				if (n.alternate !== r) throw Error(s(190));
			}
			if (n.tag !== 3) throw Error(s(188));
			return n.stateNode.current === n ? e : t;
		}
		function m(e) {
			var t = e.tag;
			if (t === 5 || t === 26 || t === 27 || t === 6) return e;
			for (e = e.child; e !== null; ) {
				if (((t = m(e)), t !== null)) return t;
				e = e.sibling;
			}
			return null;
		}
		var h = Object.assign,
			g = Symbol.for(`react.element`),
			_ = Symbol.for(`react.transitional.element`),
			v = Symbol.for(`react.portal`),
			y = Symbol.for(`react.fragment`),
			b = Symbol.for(`react.strict_mode`),
			x = Symbol.for(`react.profiler`),
			ee = Symbol.for(`react.consumer`),
			S = Symbol.for(`react.context`),
			C = Symbol.for(`react.forward_ref`),
			te = Symbol.for(`react.suspense`),
			ne = Symbol.for(`react.suspense_list`),
			re = Symbol.for(`react.memo`),
			w = Symbol.for(`react.lazy`),
			ie = Symbol.for(`react.activity`),
			ae = Symbol.for(`react.memo_cache_sentinel`),
			oe = Symbol.iterator;
		function se(e) {
			return typeof e != `object` || !e
				? null
				: ((e = (oe && e[oe]) || e[`@@iterator`]), typeof e == `function` ? e : null);
		}
		var ce = Symbol.for(`react.client.reference`);
		function le(e) {
			if (e == null) return null;
			if (typeof e == `function`) return e.$$typeof === ce ? null : e.displayName || e.name || null;
			if (typeof e == `string`) return e;
			switch (e) {
				case y:
					return `Fragment`;
				case x:
					return `Profiler`;
				case b:
					return `StrictMode`;
				case te:
					return `Suspense`;
				case ne:
					return `SuspenseList`;
				case ie:
					return `Activity`;
			}
			if (typeof e == `object`)
				switch (e.$$typeof) {
					case v:
						return `Portal`;
					case S:
						return e.displayName || `Context`;
					case ee:
						return (e._context.displayName || `Context`) + `.Consumer`;
					case C:
						var t = e.render;
						return (
							(e = e.displayName),
							(e ||=
								((e = t.displayName || t.name || ``),
								e === `` ? `ForwardRef` : `ForwardRef(` + e + `)`)),
							e
						);
					case re:
						return (t = e.displayName || null), t === null ? le(e.type) || `Memo` : t;
					case w:
						(t = e._payload), (e = e._init);
						try {
							return le(e(t));
						} catch {}
				}
			return null;
		}
		var ue = Array.isArray,
			T = r.__CLIENT_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE,
			E = a.__DOM_INTERNALS_DO_NOT_USE_OR_WARN_USERS_THEY_CANNOT_UPGRADE,
			de = { pending: !1, data: null, method: null, action: null },
			fe = [],
			pe = -1;
		function D(e) {
			return { current: e };
		}
		function O(e) {
			0 > pe || ((e.current = fe[pe]), (fe[pe] = null), pe--);
		}
		function k(e, t) {
			pe++, (fe[pe] = e.current), (e.current = t);
		}
		var me = D(null),
			he = D(null),
			ge = D(null),
			_e = D(null);
		function ve(e, t) {
			switch ((k(ge, t), k(he, e), k(me, null), t.nodeType)) {
				case 9:
				case 11:
					e = (e = t.documentElement) && (e = e.namespaceURI) ? Vd(e) : 0;
					break;
				default:
					if (((e = t.tagName), (t = t.namespaceURI))) (t = Vd(t)), (e = Hd(t, e));
					else
						switch (e) {
							case `svg`:
								e = 1;
								break;
							case `math`:
								e = 2;
								break;
							default:
								e = 0;
						}
			}
			O(me), k(me, e);
		}
		function ye() {
			O(me), O(he), O(ge);
		}
		function be(e) {
			e.memoizedState !== null && k(_e, e);
			var t = me.current,
				n = Hd(t, e.type);
			t !== n && (k(he, e), k(me, n));
		}
		function xe(e) {
			he.current === e && (O(me), O(he)), _e.current === e && (O(_e), (Qf._currentValue = de));
		}
		var Se, Ce;
		function we(e) {
			if (Se === void 0)
				try {
					throw Error();
				} catch (e) {
					var t = e.stack.trim().match(/\n( *(at )?)/);
					(Se = (t && t[1]) || ``),
						(Ce =
							-1 <
							e.stack.indexOf(`
    at`)
								? ` (<anonymous>)`
								: -1 < e.stack.indexOf(`@`)
									? `@unknown:0:0`
									: ``);
				}
			return (
				`
` +
				Se +
				e +
				Ce
			);
		}
		var Te = !1;
		function Ee(e, t) {
			if (!e || Te) return ``;
			Te = !0;
			var n = Error.prepareStackTrace;
			Error.prepareStackTrace = void 0;
			try {
				var r = {
					DetermineComponentFrameRoot: function () {
						try {
							if (t) {
								var n = function () {
									throw Error();
								};
								if (
									(Object.defineProperty(n.prototype, "props", {
										set: function () {
											throw Error();
										},
									}),
									typeof Reflect == `object` && Reflect.construct)
								) {
									try {
										Reflect.construct(n, []);
									} catch (e) {
										var r = e;
									}
									Reflect.construct(e, [], n);
								} else {
									try {
										n.call();
									} catch (e) {
										r = e;
									}
									e.call(n.prototype);
								}
							} else {
								try {
									throw Error();
								} catch (e) {
									r = e;
								}
								(n = e()) && typeof n.catch == `function` && n.catch(function () {});
							}
						} catch (e) {
							if (e && r && typeof e.stack == `string`) return [e.stack, r.stack];
						}
						return [null, null];
					},
				};
				r.DetermineComponentFrameRoot.displayName = `DetermineComponentFrameRoot`;
				var i = Object.getOwnPropertyDescriptor(r.DetermineComponentFrameRoot, `name`);
				i &&
					i.configurable &&
					Object.defineProperty(r.DetermineComponentFrameRoot, "name", {
						value: `DetermineComponentFrameRoot`,
					});
				var a = r.DetermineComponentFrameRoot(),
					o = a[0],
					s = a[1];
				if (o && s) {
					var c = o.split(`
`),
						l = s.split(`
`);
					for (i = r = 0; r < c.length && !c[r].includes(`DetermineComponentFrameRoot`); ) r++;
					for (; i < l.length && !l[i].includes(`DetermineComponentFrameRoot`); ) i++;
					if (r === c.length || i === l.length)
						for (r = c.length - 1, i = l.length - 1; 1 <= r && 0 <= i && c[r] !== l[i]; ) i--;
					for (; 1 <= r && 0 <= i; r--, i--)
						if (c[r] !== l[i]) {
							if (r !== 1 || i !== 1)
								do
									if ((r--, i--, 0 > i || c[r] !== l[i])) {
										var u =
											`
` + c[r].replace(` at new `, ` at `);
										return (
											e.displayName &&
												u.includes(`<anonymous>`) &&
												(u = u.replace(`<anonymous>`, e.displayName)),
											u
										);
									}
								while (1 <= r && 0 <= i);
							break;
						}
				}
			} finally {
				(Te = !1), (Error.prepareStackTrace = n);
			}
			return (n = e ? e.displayName || e.name : ``) ? we(n) : ``;
		}
		function De(e, t) {
			switch (e.tag) {
				case 26:
				case 27:
				case 5:
					return we(e.type);
				case 16:
					return we(`Lazy`);
				case 13:
					return e.child !== t && t !== null ? we(`Suspense Fallback`) : we(`Suspense`);
				case 19:
					return we(`SuspenseList`);
				case 0:
				case 15:
					return Ee(e.type, !1);
				case 11:
					return Ee(e.type.render, !1);
				case 1:
					return Ee(e.type, !0);
				case 31:
					return we(`Activity`);
				default:
					return ``;
			}
		}
		function A(e) {
			try {
				var t = ``,
					n = null;
				do (t += De(e, n)), (n = e), (e = e.return);
				while (e);
				return t;
			} catch (e) {
				return (
					`
Error generating stack: ` +
					e.message +
					`
` +
					e.stack
				);
			}
		}
		var Oe = Object.prototype.hasOwnProperty,
			ke = t.unstable_scheduleCallback,
			Ae = t.unstable_cancelCallback,
			je = t.unstable_shouldYield,
			Me = t.unstable_requestPaint,
			Ne = t.unstable_now,
			Pe = t.unstable_getCurrentPriorityLevel,
			Fe = t.unstable_ImmediatePriority,
			Ie = t.unstable_UserBlockingPriority,
			Le = t.unstable_NormalPriority,
			Re = t.unstable_LowPriority,
			ze = t.unstable_IdlePriority,
			Be = t.log,
			Ve = t.unstable_setDisableYieldValue,
			He = null,
			Ue = null;
		function We(e) {
			if ((typeof Be == `function` && Ve(e), Ue && typeof Ue.setStrictMode == `function`))
				try {
					Ue.setStrictMode(He, e);
				} catch {}
		}
		var Ge = Math.clz32 ? Math.clz32 : Je,
			Ke = Math.log,
			qe = Math.LN2;
		function Je(e) {
			return (e >>>= 0), e === 0 ? 32 : (31 - ((Ke(e) / qe) | 0)) | 0;
		}
		var Ye = 256,
			Xe = 262144,
			Ze = 4194304;
		function Qe(e) {
			var t = e & 42;
			if (t !== 0) return t;
			switch (e & -e) {
				case 1:
					return 1;
				case 2:
					return 2;
				case 4:
					return 4;
				case 8:
					return 8;
				case 16:
					return 16;
				case 32:
					return 32;
				case 64:
					return 64;
				case 128:
					return 128;
				case 256:
				case 512:
				case 1024:
				case 2048:
				case 4096:
				case 8192:
				case 16384:
				case 32768:
				case 65536:
				case 131072:
					return e & 261888;
				case 262144:
				case 524288:
				case 1048576:
				case 2097152:
					return e & 3932160;
				case 4194304:
				case 8388608:
				case 16777216:
				case 33554432:
					return e & 62914560;
				case 67108864:
					return 67108864;
				case 134217728:
					return 134217728;
				case 268435456:
					return 268435456;
				case 536870912:
					return 536870912;
				case 1073741824:
					return 0;
				default:
					return e;
			}
		}
		function $e(e, t, n) {
			var r = e.pendingLanes;
			if (r === 0) return 0;
			var i = 0,
				a = e.suspendedLanes,
				o = e.pingedLanes;
			e = e.warmLanes;
			var s = r & 134217727;
			return (
				s === 0
					? ((s = r & ~a),
						s === 0
							? o === 0
								? n || ((n = r & ~e), n !== 0 && (i = Qe(n)))
								: (i = Qe(o))
							: (i = Qe(s)))
					: ((r = s & ~a),
						r === 0
							? ((o &= s), o === 0 ? n || ((n = s & ~e), n !== 0 && (i = Qe(n))) : (i = Qe(o)))
							: (i = Qe(r))),
				i === 0
					? 0
					: t !== 0 &&
							t !== i &&
							(t & a) === 0 &&
							((a = i & -i), (n = t & -t), a >= n || (a === 32 && n & 4194048))
						? t
						: i
			);
		}
		function et(e, t) {
			return (e.pendingLanes & ~(e.suspendedLanes & ~e.pingedLanes) & t) === 0;
		}
		function tt(e, t) {
			switch (e) {
				case 1:
				case 2:
				case 4:
				case 8:
				case 64:
					return t + 250;
				case 16:
				case 32:
				case 128:
				case 256:
				case 512:
				case 1024:
				case 2048:
				case 4096:
				case 8192:
				case 16384:
				case 32768:
				case 65536:
				case 131072:
				case 262144:
				case 524288:
				case 1048576:
				case 2097152:
					return t + 5e3;
				case 4194304:
				case 8388608:
				case 16777216:
				case 33554432:
					return -1;
				case 67108864:
				case 134217728:
				case 268435456:
				case 536870912:
				case 1073741824:
					return -1;
				default:
					return -1;
			}
		}
		function nt() {
			var e = Ze;
			return (Ze <<= 1), !(Ze & 62914560) && (Ze = 4194304), e;
		}
		function rt(e) {
			for (var t = [], n = 0; 31 > n; n++) t.push(e);
			return t;
		}
		function it(e, t) {
			(e.pendingLanes |= t),
				t !== 268435456 && ((e.suspendedLanes = 0), (e.pingedLanes = 0), (e.warmLanes = 0));
		}
		function at(e, t, n, r, i, a) {
			var o = e.pendingLanes;
			(e.pendingLanes = n),
				(e.suspendedLanes = 0),
				(e.pingedLanes = 0),
				(e.warmLanes = 0),
				(e.expiredLanes &= n),
				(e.entangledLanes &= n),
				(e.errorRecoveryDisabledLanes &= n),
				(e.shellSuspendCounter = 0);
			var s = e.entanglements,
				c = e.expirationTimes,
				l = e.hiddenUpdates;
			for (n = o & ~n; 0 < n; ) {
				var u = 31 - Ge(n),
					d = 1 << u;
				(s[u] = 0), (c[u] = -1);
				var f = l[u];
				if (f !== null)
					for (l[u] = null, u = 0; u < f.length; u++) {
						var p = f[u];
						p !== null && (p.lane &= -536870913);
					}
				n &= ~d;
			}
			r !== 0 && ot(e, r, 0),
				a !== 0 && i === 0 && e.tag !== 0 && (e.suspendedLanes |= a & ~(o & ~t));
		}
		function ot(e, t, n) {
			(e.pendingLanes |= t), (e.suspendedLanes &= ~t);
			var r = 31 - Ge(t);
			(e.entangledLanes |= t),
				(e.entanglements[r] = e.entanglements[r] | 1073741824 | (n & 261930));
		}
		function st(e, t) {
			var n = (e.entangledLanes |= t);
			for (e = e.entanglements; n; ) {
				var r = 31 - Ge(n),
					i = 1 << r;
				(i & t) | (e[r] & t) && (e[r] |= t), (n &= ~i);
			}
		}
		function ct(e, t) {
			var n = t & -t;
			return (n = n & 42 ? 1 : lt(n)), (n & (e.suspendedLanes | t)) === 0 ? n : 0;
		}
		function lt(e) {
			switch (e) {
				case 2:
					e = 1;
					break;
				case 8:
					e = 4;
					break;
				case 32:
					e = 16;
					break;
				case 256:
				case 512:
				case 1024:
				case 2048:
				case 4096:
				case 8192:
				case 16384:
				case 32768:
				case 65536:
				case 131072:
				case 262144:
				case 524288:
				case 1048576:
				case 2097152:
				case 4194304:
				case 8388608:
				case 16777216:
				case 33554432:
					e = 128;
					break;
				case 268435456:
					e = 134217728;
					break;
				default:
					e = 0;
			}
			return e;
		}
		function ut(e) {
			return (e &= -e), 2 < e ? (8 < e ? (e & 134217727 ? 32 : 268435456) : 8) : 2;
		}
		function dt() {
			var e = E.p;
			return e === 0 ? ((e = window.event), e === void 0 ? 32 : mp(e.type)) : e;
		}
		function ft(e, t) {
			var n = E.p;
			try {
				return (E.p = e), t();
			} finally {
				E.p = n;
			}
		}
		var pt = Math.random().toString(36).slice(2),
			mt = `__reactFiber$` + pt,
			ht = `__reactProps$` + pt,
			gt = `__reactContainer$` + pt,
			_t = `__reactEvents$` + pt,
			vt = `__reactListeners$` + pt,
			yt = `__reactHandles$` + pt,
			bt = `__reactResources$` + pt,
			xt = `__reactMarker$` + pt;
		function St(e) {
			delete e[mt], delete e[ht], delete e[_t], delete e[vt], delete e[yt];
		}
		function Ct(e) {
			var t = e[mt];
			if (t) return t;
			for (var n = e.parentNode; n; ) {
				if ((t = n[gt] || n[mt])) {
					if (((n = t.alternate), t.child !== null || (n !== null && n.child !== null)))
						for (e = df(e); e !== null; ) {
							if ((n = e[mt])) return n;
							e = df(e);
						}
					return t;
				}
				(e = n), (n = e.parentNode);
			}
			return null;
		}
		function wt(e) {
			if ((e = e[mt] || e[gt])) {
				var t = e.tag;
				if (t === 5 || t === 6 || t === 13 || t === 31 || t === 26 || t === 27 || t === 3) return e;
			}
			return null;
		}
		function Tt(e) {
			var t = e.tag;
			if (t === 5 || t === 26 || t === 27 || t === 6) return e.stateNode;
			throw Error(s(33));
		}
		function Et(e) {
			var t = e[bt];
			return (t ||= e[bt] = { hoistableStyles: new Map(), hoistableScripts: new Map() }), t;
		}
		function Dt(e) {
			e[xt] = !0;
		}
		var Ot = new Set(),
			kt = {};
		function At(e, t) {
			jt(e, t), jt(e + `Capture`, t);
		}
		function jt(e, t) {
			for (kt[e] = t, e = 0; e < t.length; e++) Ot.add(t[e]);
		}
		var Mt = RegExp(
				`^[:A-Z_a-z\\u00C0-\\u00D6\\u00D8-\\u00F6\\u00F8-\\u02FF\\u0370-\\u037D\\u037F-\\u1FFF\\u200C-\\u200D\\u2070-\\u218F\\u2C00-\\u2FEF\\u3001-\\uD7FF\\uF900-\\uFDCF\\uFDF0-\\uFFFD][:A-Z_a-z\\u00C0-\\u00D6\\u00D8-\\u00F6\\u00F8-\\u02FF\\u0370-\\u037D\\u037F-\\u1FFF\\u200C-\\u200D\\u2070-\\u218F\\u2C00-\\u2FEF\\u3001-\\uD7FF\\uF900-\\uFDCF\\uFDF0-\\uFFFD\\-.0-9\\u00B7\\u0300-\\u036F\\u203F-\\u2040]*$`,
			),
			Nt = {},
			Pt = {};
		function Ft(e) {
			return Oe.call(Pt, e)
				? !0
				: Oe.call(Nt, e)
					? !1
					: Mt.test(e)
						? (Pt[e] = !0)
						: ((Nt[e] = !0), !1);
		}
		function It(e, t, n) {
			if (Ft(t))
				if (n === null) e.removeAttribute(t);
				else {
					switch (typeof n) {
						case `undefined`:
						case `function`:
						case `symbol`:
							e.removeAttribute(t);
							return;
						case `boolean`:
							var r = t.toLowerCase().slice(0, 5);
							if (r !== `data-` && r !== `aria-`) {
								e.removeAttribute(t);
								return;
							}
					}
					e.setAttribute(t, `` + n);
				}
		}
		function Lt(e, t, n) {
			if (n === null) e.removeAttribute(t);
			else {
				switch (typeof n) {
					case `undefined`:
					case `function`:
					case `symbol`:
					case `boolean`:
						e.removeAttribute(t);
						return;
				}
				e.setAttribute(t, `` + n);
			}
		}
		function Rt(e, t, n, r) {
			if (r === null) e.removeAttribute(n);
			else {
				switch (typeof r) {
					case `undefined`:
					case `function`:
					case `symbol`:
					case `boolean`:
						e.removeAttribute(n);
						return;
				}
				e.setAttributeNS(t, n, `` + r);
			}
		}
		function zt(e) {
			switch (typeof e) {
				case `bigint`:
				case `boolean`:
				case `number`:
				case `string`:
				case `undefined`:
					return e;
				case `object`:
					return e;
				default:
					return ``;
			}
		}
		function Bt(e) {
			var t = e.type;
			return (e = e.nodeName) && e.toLowerCase() === `input` && (t === `checkbox` || t === `radio`);
		}
		function Vt(e, t, n) {
			var r = Object.getOwnPropertyDescriptor(e.constructor.prototype, t);
			if (
				!e.hasOwnProperty(t) &&
				r !== void 0 &&
				typeof r.get == `function` &&
				typeof r.set == `function`
			) {
				var i = r.get,
					a = r.set;
				return (
					Object.defineProperty(e, t, {
						configurable: !0,
						get: function () {
							return i.call(this);
						},
						set: function (e) {
							(n = `` + e), a.call(this, e);
						},
					}),
					Object.defineProperty(e, t, { enumerable: r.enumerable }),
					{
						getValue: function () {
							return n;
						},
						setValue: function (e) {
							n = `` + e;
						},
						stopTracking: function () {
							(e._valueTracker = null), delete e[t];
						},
					}
				);
			}
		}
		function Ht(e) {
			if (!e._valueTracker) {
				var t = Bt(e) ? `checked` : `value`;
				e._valueTracker = Vt(e, t, `` + e[t]);
			}
		}
		function Ut(e) {
			if (!e) return !1;
			var t = e._valueTracker;
			if (!t) return !0;
			var n = t.getValue(),
				r = ``;
			return (
				e && (r = Bt(e) ? (e.checked ? `true` : `false`) : e.value),
				(e = r),
				e === n ? !1 : (t.setValue(e), !0)
			);
		}
		function Wt(e) {
			if (((e ||= typeof document < `u` ? document : void 0), e === void 0)) return null;
			try {
				return e.activeElement || e.body;
			} catch {
				return e.body;
			}
		}
		var Gt = /[\n"\\]/g;
		function Kt(e) {
			return e.replace(Gt, function (e) {
				return `\\` + e.charCodeAt(0).toString(16) + ` `;
			});
		}
		function qt(e, t, n, r, i, a, o, s) {
			(e.name = ``),
				o != null && typeof o != `function` && typeof o != `symbol` && typeof o != `boolean`
					? (e.type = o)
					: e.removeAttribute(`type`),
				t == null
					? (o !== `submit` && o !== `reset`) || e.removeAttribute(`value`)
					: o === `number`
						? ((t === 0 && e.value === ``) || e.value != t) && (e.value = `` + zt(t))
						: e.value !== `` + zt(t) && (e.value = `` + zt(t)),
				t == null
					? n == null
						? r != null && e.removeAttribute(`value`)
						: Yt(e, o, zt(n))
					: Yt(e, o, zt(t)),
				i == null && a != null && (e.defaultChecked = !!a),
				i != null && (e.checked = i && typeof i != `function` && typeof i != `symbol`),
				s != null && typeof s != `function` && typeof s != `symbol` && typeof s != `boolean`
					? (e.name = `` + zt(s))
					: e.removeAttribute(`name`);
		}
		function Jt(e, t, n, r, i, a, o, s) {
			if (
				(a != null &&
					typeof a != `function` &&
					typeof a != `symbol` &&
					typeof a != `boolean` &&
					(e.type = a),
				t != null || n != null)
			) {
				if (!((a !== `submit` && a !== `reset`) || t != null)) {
					Ht(e);
					return;
				}
				(n = n == null ? `` : `` + zt(n)),
					(t = t == null ? n : `` + zt(t)),
					s || t === e.value || (e.value = t),
					(e.defaultValue = t);
			}
			(r ??= i),
				(r = typeof r != `function` && typeof r != `symbol` && !!r),
				(e.checked = s ? e.checked : !!r),
				(e.defaultChecked = !!r),
				o != null &&
					typeof o != `function` &&
					typeof o != `symbol` &&
					typeof o != `boolean` &&
					(e.name = o),
				Ht(e);
		}
		function Yt(e, t, n) {
			(t === `number` && Wt(e.ownerDocument) === e) ||
				e.defaultValue === `` + n ||
				(e.defaultValue = `` + n);
		}
		function Xt(e, t, n, r) {
			if (((e = e.options), t)) {
				t = {};
				for (var i = 0; i < n.length; i++) t[`$` + n[i]] = !0;
				for (n = 0; n < e.length; n++)
					(i = t.hasOwnProperty(`$` + e[n].value)),
						e[n].selected !== i && (e[n].selected = i),
						i && r && (e[n].defaultSelected = !0);
			} else {
				for (n = `` + zt(n), t = null, i = 0; i < e.length; i++) {
					if (e[i].value === n) {
						(e[i].selected = !0), r && (e[i].defaultSelected = !0);
						return;
					}
					t !== null || e[i].disabled || (t = e[i]);
				}
				t !== null && (t.selected = !0);
			}
		}
		function Zt(e, t, n) {
			if (t != null && ((t = `` + zt(t)), t !== e.value && (e.value = t), n == null)) {
				e.defaultValue !== t && (e.defaultValue = t);
				return;
			}
			e.defaultValue = n == null ? `` : `` + zt(n);
		}
		function Qt(e, t, n, r) {
			if (t == null) {
				if (r != null) {
					if (n != null) throw Error(s(92));
					if (ue(r)) {
						if (1 < r.length) throw Error(s(93));
						r = r[0];
					}
					n = r;
				}
				(n ??= ``), (t = n);
			}
			(n = zt(t)),
				(e.defaultValue = n),
				(r = e.textContent),
				r === n && r !== `` && r !== null && (e.value = r),
				Ht(e);
		}
		function j(e, t) {
			if (t) {
				var n = e.firstChild;
				if (n && n === e.lastChild && n.nodeType === 3) {
					n.nodeValue = t;
					return;
				}
			}
			e.textContent = t;
		}
		var $t = new Set(
			`animationIterationCount aspectRatio borderImageOutset borderImageSlice borderImageWidth boxFlex boxFlexGroup boxOrdinalGroup columnCount columns flex flexGrow flexPositive flexShrink flexNegative flexOrder gridArea gridRow gridRowEnd gridRowSpan gridRowStart gridColumn gridColumnEnd gridColumnSpan gridColumnStart fontWeight lineClamp lineHeight opacity order orphans scale tabSize widows zIndex zoom fillOpacity floodOpacity stopOpacity strokeDasharray strokeDashoffset strokeMiterlimit strokeOpacity strokeWidth MozAnimationIterationCount MozBoxFlex MozBoxFlexGroup MozLineClamp msAnimationIterationCount msFlex msZoom msFlexGrow msFlexNegative msFlexOrder msFlexPositive msFlexShrink msGridColumn msGridColumnSpan msGridRow msGridRowSpan WebkitAnimationIterationCount WebkitBoxFlex WebKitBoxFlexGroup WebkitBoxOrdinalGroup WebkitColumnCount WebkitColumns WebkitFlex WebkitFlexGrow WebkitFlexPositive WebkitFlexShrink WebkitLineClamp`.split(
				` `,
			),
		);
		function en(e, t, n) {
			var r = t.indexOf(`--`) === 0;
			n == null || typeof n == `boolean` || n === ``
				? r
					? e.setProperty(t, ``)
					: t === `float`
						? (e.cssFloat = ``)
						: (e[t] = ``)
				: r
					? e.setProperty(t, n)
					: typeof n != `number` || n === 0 || $t.has(t)
						? t === `float`
							? (e.cssFloat = n)
							: (e[t] = (`` + n).trim())
						: (e[t] = n + `px`);
		}
		function tn(e, t, n) {
			if (t != null && typeof t != `object`) throw Error(s(62));
			if (((e = e.style), n != null)) {
				for (var r in n)
					!n.hasOwnProperty(r) ||
						(t != null && t.hasOwnProperty(r)) ||
						(r.indexOf(`--`) === 0
							? e.setProperty(r, ``)
							: r === `float`
								? (e.cssFloat = ``)
								: (e[r] = ``));
				for (var i in t) (r = t[i]), t.hasOwnProperty(i) && n[i] !== r && en(e, i, r);
			} else for (var a in t) t.hasOwnProperty(a) && en(e, a, t[a]);
		}
		function nn(e) {
			if (e.indexOf(`-`) === -1) return !1;
			switch (e) {
				case `annotation-xml`:
				case `color-profile`:
				case `font-face`:
				case `font-face-src`:
				case `font-face-uri`:
				case `font-face-format`:
				case `font-face-name`:
				case `missing-glyph`:
					return !1;
				default:
					return !0;
			}
		}
		var rn = new Map([
				[`acceptCharset`, `accept-charset`],
				[`htmlFor`, `for`],
				[`httpEquiv`, `http-equiv`],
				[`crossOrigin`, `crossorigin`],
				[`accentHeight`, `accent-height`],
				[`alignmentBaseline`, `alignment-baseline`],
				[`arabicForm`, `arabic-form`],
				[`baselineShift`, `baseline-shift`],
				[`capHeight`, `cap-height`],
				[`clipPath`, `clip-path`],
				[`clipRule`, `clip-rule`],
				[`colorInterpolation`, `color-interpolation`],
				[`colorInterpolationFilters`, `color-interpolation-filters`],
				[`colorProfile`, `color-profile`],
				[`colorRendering`, `color-rendering`],
				[`dominantBaseline`, `dominant-baseline`],
				[`enableBackground`, `enable-background`],
				[`fillOpacity`, `fill-opacity`],
				[`fillRule`, `fill-rule`],
				[`floodColor`, `flood-color`],
				[`floodOpacity`, `flood-opacity`],
				[`fontFamily`, `font-family`],
				[`fontSize`, `font-size`],
				[`fontSizeAdjust`, `font-size-adjust`],
				[`fontStretch`, `font-stretch`],
				[`fontStyle`, `font-style`],
				[`fontVariant`, `font-variant`],
				[`fontWeight`, `font-weight`],
				[`glyphName`, `glyph-name`],
				[`glyphOrientationHorizontal`, `glyph-orientation-horizontal`],
				[`glyphOrientationVertical`, `glyph-orientation-vertical`],
				[`horizAdvX`, `horiz-adv-x`],
				[`horizOriginX`, `horiz-origin-x`],
				[`imageRendering`, `image-rendering`],
				[`letterSpacing`, `letter-spacing`],
				[`lightingColor`, `lighting-color`],
				[`markerEnd`, `marker-end`],
				[`markerMid`, `marker-mid`],
				[`markerStart`, `marker-start`],
				[`overlinePosition`, `overline-position`],
				[`overlineThickness`, `overline-thickness`],
				[`paintOrder`, `paint-order`],
				[`panose-1`, `panose-1`],
				[`pointerEvents`, `pointer-events`],
				[`renderingIntent`, `rendering-intent`],
				[`shapeRendering`, `shape-rendering`],
				[`stopColor`, `stop-color`],
				[`stopOpacity`, `stop-opacity`],
				[`strikethroughPosition`, `strikethrough-position`],
				[`strikethroughThickness`, `strikethrough-thickness`],
				[`strokeDasharray`, `stroke-dasharray`],
				[`strokeDashoffset`, `stroke-dashoffset`],
				[`strokeLinecap`, `stroke-linecap`],
				[`strokeLinejoin`, `stroke-linejoin`],
				[`strokeMiterlimit`, `stroke-miterlimit`],
				[`strokeOpacity`, `stroke-opacity`],
				[`strokeWidth`, `stroke-width`],
				[`textAnchor`, `text-anchor`],
				[`textDecoration`, `text-decoration`],
				[`textRendering`, `text-rendering`],
				[`transformOrigin`, `transform-origin`],
				[`underlinePosition`, `underline-position`],
				[`underlineThickness`, `underline-thickness`],
				[`unicodeBidi`, `unicode-bidi`],
				[`unicodeRange`, `unicode-range`],
				[`unitsPerEm`, `units-per-em`],
				[`vAlphabetic`, `v-alphabetic`],
				[`vHanging`, `v-hanging`],
				[`vIdeographic`, `v-ideographic`],
				[`vMathematical`, `v-mathematical`],
				[`vectorEffect`, `vector-effect`],
				[`vertAdvY`, `vert-adv-y`],
				[`vertOriginX`, `vert-origin-x`],
				[`vertOriginY`, `vert-origin-y`],
				[`wordSpacing`, `word-spacing`],
				[`writingMode`, `writing-mode`],
				[`xmlnsXlink`, `xmlns:xlink`],
				[`xHeight`, `x-height`],
			]),
			an =
				/^[\u0000-\u001F ]*j[\r\n\t]*a[\r\n\t]*v[\r\n\t]*a[\r\n\t]*s[\r\n\t]*c[\r\n\t]*r[\r\n\t]*i[\r\n\t]*p[\r\n\t]*t[\r\n\t]*:/i;
		function on(e) {
			return an.test(`` + e)
				? `javascript:throw new Error('React has blocked a javascript: URL as a security precaution.')`
				: e;
		}
		function sn() {}
		var cn = null;
		function ln(e) {
			return (
				(e = e.target || e.srcElement || window),
				e.correspondingUseElement && (e = e.correspondingUseElement),
				e.nodeType === 3 ? e.parentNode : e
			);
		}
		var un = null,
			dn = null;
		function fn(e) {
			var t = wt(e);
			if (t && (e = t.stateNode)) {
				var n = e[ht] || null;
				a: switch (((e = t.stateNode), t.type)) {
					case `input`:
						if (
							(qt(
								e,
								n.value,
								n.defaultValue,
								n.defaultValue,
								n.checked,
								n.defaultChecked,
								n.type,
								n.name,
							),
							(t = n.name),
							n.type === `radio` && t != null)
						) {
							for (n = e; n.parentNode; ) n = n.parentNode;
							for (
								n = n.querySelectorAll(`input[name="` + Kt(`` + t) + `"][type="radio"]`), t = 0;
								t < n.length;
								t++
							) {
								var r = n[t];
								if (r !== e && r.form === e.form) {
									var i = r[ht] || null;
									if (!i) throw Error(s(90));
									qt(
										r,
										i.value,
										i.defaultValue,
										i.defaultValue,
										i.checked,
										i.defaultChecked,
										i.type,
										i.name,
									);
								}
							}
							for (t = 0; t < n.length; t++) (r = n[t]), r.form === e.form && Ut(r);
						}
						break a;
					case `textarea`:
						Zt(e, n.value, n.defaultValue);
						break a;
					case `select`:
						(t = n.value), t != null && Xt(e, !!n.multiple, t, !1);
				}
			}
		}
		var pn = !1;
		function mn(e, t, n) {
			if (pn) return e(t, n);
			pn = !0;
			try {
				return e(t);
			} finally {
				if (
					((pn = !1),
					(un !== null || dn !== null) &&
						(bu(), un && ((t = un), (e = dn), (dn = un = null), fn(t), e)))
				)
					for (t = 0; t < e.length; t++) fn(e[t]);
			}
		}
		function hn(e, t) {
			var n = e.stateNode;
			if (n === null) return null;
			var r = n[ht] || null;
			if (r === null) return null;
			n = r[t];
			a: switch (t) {
				case `onClick`:
				case `onClickCapture`:
				case `onDoubleClick`:
				case `onDoubleClickCapture`:
				case `onMouseDown`:
				case `onMouseDownCapture`:
				case `onMouseMove`:
				case `onMouseMoveCapture`:
				case `onMouseUp`:
				case `onMouseUpCapture`:
				case `onMouseEnter`:
					(r = !r.disabled) ||
						((e = e.type),
						(r = !(e === `button` || e === `input` || e === `select` || e === `textarea`))),
						(e = !r);
					break a;
				default:
					e = !1;
			}
			if (e) return null;
			if (n && typeof n != `function`) throw Error(s(231, t, typeof n));
			return n;
		}
		var gn = !(
				typeof window > `u` ||
				window.document === void 0 ||
				window.document.createElement === void 0
			),
			_n = !1;
		if (gn)
			try {
				var vn = {};
				Object.defineProperty(vn, "passive", {
					get: function () {
						_n = !0;
					},
				}),
					window.addEventListener(`test`, vn, vn),
					window.removeEventListener(`test`, vn, vn);
			} catch {
				_n = !1;
			}
		var yn = null,
			bn = null,
			xn = null;
		function Sn() {
			if (xn) return xn;
			var e,
				t = bn,
				n = t.length,
				r,
				i = `value` in yn ? yn.value : yn.textContent,
				a = i.length;
			for (e = 0; e < n && t[e] === i[e]; e++);
			var o = n - e;
			for (r = 1; r <= o && t[n - r] === i[a - r]; r++);
			return (xn = i.slice(e, 1 < r ? 1 - r : void 0));
		}
		function Cn(e) {
			var t = e.keyCode;
			return (
				`charCode` in e ? ((e = e.charCode), e === 0 && t === 13 && (e = 13)) : (e = t),
				e === 10 && (e = 13),
				32 <= e || e === 13 ? e : 0
			);
		}
		function wn() {
			return !0;
		}
		function Tn() {
			return !1;
		}
		function En(e) {
			function t(t, n, r, i, a) {
				for (var o in ((this._reactName = t),
				(this._targetInst = r),
				(this.type = n),
				(this.nativeEvent = i),
				(this.target = a),
				(this.currentTarget = null),
				e))
					e.hasOwnProperty(o) && ((t = e[o]), (this[o] = t ? t(i) : i[o]));
				return (
					(this.isDefaultPrevented = (
						i.defaultPrevented == null
							? !1 === i.returnValue
							: i.defaultPrevented
					)
						? wn
						: Tn),
					(this.isPropagationStopped = Tn),
					this
				);
			}
			return (
				h(t.prototype, {
					preventDefault: function () {
						this.defaultPrevented = !0;
						var e = this.nativeEvent;
						e &&
							(e.preventDefault
								? e.preventDefault()
								: typeof e.returnValue != `unknown` && (e.returnValue = !1),
							(this.isDefaultPrevented = wn));
					},
					stopPropagation: function () {
						var e = this.nativeEvent;
						e &&
							(e.stopPropagation
								? e.stopPropagation()
								: typeof e.cancelBubble != `unknown` && (e.cancelBubble = !0),
							(this.isPropagationStopped = wn));
					},
					persist: function () {},
					isPersistent: wn,
				}),
				t
			);
		}
		var Dn = {
				eventPhase: 0,
				bubbles: 0,
				cancelable: 0,
				timeStamp: function (e) {
					return e.timeStamp || Date.now();
				},
				defaultPrevented: 0,
				isTrusted: 0,
			},
			On = En(Dn),
			kn = h({}, Dn, { view: 0, detail: 0 }),
			An = En(kn),
			jn,
			Mn,
			Nn,
			Pn = h({}, kn, {
				screenX: 0,
				screenY: 0,
				clientX: 0,
				clientY: 0,
				pageX: 0,
				pageY: 0,
				ctrlKey: 0,
				shiftKey: 0,
				altKey: 0,
				metaKey: 0,
				getModifierState: Gn,
				button: 0,
				buttons: 0,
				relatedTarget: function (e) {
					return e.relatedTarget === void 0
						? e.fromElement === e.srcElement
							? e.toElement
							: e.fromElement
						: e.relatedTarget;
				},
				movementX: function (e) {
					return `movementX` in e
						? e.movementX
						: (e !== Nn &&
								(Nn && e.type === `mousemove`
									? ((jn = e.screenX - Nn.screenX), (Mn = e.screenY - Nn.screenY))
									: (Mn = jn = 0),
								(Nn = e)),
							jn);
				},
				movementY: function (e) {
					return `movementY` in e ? e.movementY : Mn;
				},
			}),
			Fn = En(Pn),
			In = En(h({}, Pn, { dataTransfer: 0 })),
			Ln = En(h({}, kn, { relatedTarget: 0 })),
			Rn = En(h({}, Dn, { animationName: 0, elapsedTime: 0, pseudoElement: 0 })),
			zn = En(
				h({}, Dn, {
					clipboardData: function (e) {
						return `clipboardData` in e ? e.clipboardData : window.clipboardData;
					},
				}),
			),
			Bn = En(h({}, Dn, { data: 0 })),
			Vn = {
				Esc: `Escape`,
				Spacebar: ` `,
				Left: `ArrowLeft`,
				Up: `ArrowUp`,
				Right: `ArrowRight`,
				Down: `ArrowDown`,
				Del: `Delete`,
				Win: `OS`,
				Menu: `ContextMenu`,
				Apps: `ContextMenu`,
				Scroll: `ScrollLock`,
				MozPrintableKey: `Unidentified`,
			},
			Hn = {
				8: `Backspace`,
				9: `Tab`,
				12: `Clear`,
				13: `Enter`,
				16: `Shift`,
				17: `Control`,
				18: `Alt`,
				19: `Pause`,
				20: `CapsLock`,
				27: `Escape`,
				32: ` `,
				33: `PageUp`,
				34: `PageDown`,
				35: `End`,
				36: `Home`,
				37: `ArrowLeft`,
				38: `ArrowUp`,
				39: `ArrowRight`,
				40: `ArrowDown`,
				45: `Insert`,
				46: `Delete`,
				112: `F1`,
				113: `F2`,
				114: `F3`,
				115: `F4`,
				116: `F5`,
				117: `F6`,
				118: `F7`,
				119: `F8`,
				120: `F9`,
				121: `F10`,
				122: `F11`,
				123: `F12`,
				144: `NumLock`,
				145: `ScrollLock`,
				224: `Meta`,
			},
			Un = { Alt: `altKey`, Control: `ctrlKey`, Meta: `metaKey`, Shift: `shiftKey` };
		function Wn(e) {
			var t = this.nativeEvent;
			return t.getModifierState ? t.getModifierState(e) : (e = Un[e]) ? !!t[e] : !1;
		}
		function Gn() {
			return Wn;
		}
		var Kn = En(
				h({}, kn, {
					key: function (e) {
						if (e.key) {
							var t = Vn[e.key] || e.key;
							if (t !== `Unidentified`) return t;
						}
						return e.type === `keypress`
							? ((e = Cn(e)), e === 13 ? `Enter` : String.fromCharCode(e))
							: e.type === `keydown` || e.type === `keyup`
								? Hn[e.keyCode] || `Unidentified`
								: ``;
					},
					code: 0,
					location: 0,
					ctrlKey: 0,
					shiftKey: 0,
					altKey: 0,
					metaKey: 0,
					repeat: 0,
					locale: 0,
					getModifierState: Gn,
					charCode: function (e) {
						return e.type === `keypress` ? Cn(e) : 0;
					},
					keyCode: function (e) {
						return e.type === `keydown` || e.type === `keyup` ? e.keyCode : 0;
					},
					which: function (e) {
						return e.type === `keypress`
							? Cn(e)
							: e.type === `keydown` || e.type === `keyup`
								? e.keyCode
								: 0;
					},
				}),
			),
			qn = En(
				h({}, Pn, {
					pointerId: 0,
					width: 0,
					height: 0,
					pressure: 0,
					tangentialPressure: 0,
					tiltX: 0,
					tiltY: 0,
					twist: 0,
					pointerType: 0,
					isPrimary: 0,
				}),
			),
			Jn = En(
				h({}, kn, {
					touches: 0,
					targetTouches: 0,
					changedTouches: 0,
					altKey: 0,
					metaKey: 0,
					ctrlKey: 0,
					shiftKey: 0,
					getModifierState: Gn,
				}),
			),
			M = En(h({}, Dn, { propertyName: 0, elapsedTime: 0, pseudoElement: 0 })),
			Yn = En(
				h({}, Pn, {
					deltaX: function (e) {
						return `deltaX` in e ? e.deltaX : `wheelDeltaX` in e ? -e.wheelDeltaX : 0;
					},
					deltaY: function (e) {
						return `deltaY` in e
							? e.deltaY
							: `wheelDeltaY` in e
								? -e.wheelDeltaY
								: `wheelDelta` in e
									? -e.wheelDelta
									: 0;
					},
					deltaZ: 0,
					deltaMode: 0,
				}),
			),
			Xn = En(h({}, Dn, { newState: 0, oldState: 0 })),
			Zn = [9, 13, 27, 32],
			Qn = gn && `CompositionEvent` in window,
			$n = null;
		gn && `documentMode` in document && ($n = document.documentMode);
		var er = gn && `TextEvent` in window && !$n,
			tr = gn && (!Qn || ($n && 8 < $n && 11 >= $n)),
			nr = ` `,
			rr = !1;
		function ir(e, t) {
			switch (e) {
				case `keyup`:
					return Zn.indexOf(t.keyCode) !== -1;
				case `keydown`:
					return t.keyCode !== 229;
				case `keypress`:
				case `mousedown`:
				case `focusout`:
					return !0;
				default:
					return !1;
			}
		}
		function ar(e) {
			return (e = e.detail), typeof e == `object` && `data` in e ? e.data : null;
		}
		var or = !1;
		function sr(e, t) {
			switch (e) {
				case `compositionend`:
					return ar(t);
				case `keypress`:
					return t.which === 32 ? ((rr = !0), nr) : null;
				case `textInput`:
					return (e = t.data), e === nr && rr ? null : e;
				default:
					return null;
			}
		}
		function cr(e, t) {
			if (or)
				return e === `compositionend` || (!Qn && ir(e, t))
					? ((e = Sn()), (xn = bn = yn = null), (or = !1), e)
					: null;
			switch (e) {
				case `paste`:
					return null;
				case `keypress`:
					if (!(t.ctrlKey || t.altKey || t.metaKey) || (t.ctrlKey && t.altKey)) {
						if (t.char && 1 < t.char.length) return t.char;
						if (t.which) return String.fromCharCode(t.which);
					}
					return null;
				case `compositionend`:
					return tr && t.locale !== `ko` ? null : t.data;
				default:
					return null;
			}
		}
		var lr = {
			color: !0,
			date: !0,
			datetime: !0,
			"datetime-local": !0,
			email: !0,
			month: !0,
			number: !0,
			password: !0,
			range: !0,
			search: !0,
			tel: !0,
			text: !0,
			time: !0,
			url: !0,
			week: !0,
		};
		function ur(e) {
			var t = e && e.nodeName && e.nodeName.toLowerCase();
			return t === `input` ? !!lr[e.type] : t === `textarea`;
		}
		function dr(e, t, n, r) {
			un ? (dn ? dn.push(r) : (dn = [r])) : (un = r),
				(t = Ed(t, `onChange`)),
				0 < t.length &&
					((n = new On(`onChange`, `change`, null, n, r)), e.push({ event: n, listeners: t }));
		}
		var fr = null,
			pr = null;
		function mr(e) {
			yd(e, 0);
		}
		function hr(e) {
			if (Ut(Tt(e))) return e;
		}
		function gr(e, t) {
			if (e === `change`) return t;
		}
		var _r = !1;
		if (gn) {
			var vr;
			if (gn) {
				var yr = `oninput` in document;
				if (!yr) {
					var br = document.createElement(`div`);
					br.setAttribute(`oninput`, `return;`), (yr = typeof br.oninput == `function`);
				}
				vr = yr;
			} else vr = !1;
			_r = vr && (!document.documentMode || 9 < document.documentMode);
		}
		function xr() {
			fr && (fr.detachEvent(`onpropertychange`, Sr), (pr = fr = null));
		}
		function Sr(e) {
			if (e.propertyName === `value` && hr(pr)) {
				var t = [];
				dr(t, pr, e, ln(e)), mn(mr, t);
			}
		}
		function Cr(e, t, n) {
			e === `focusin`
				? (xr(), (fr = t), (pr = n), fr.attachEvent(`onpropertychange`, Sr))
				: e === `focusout` && xr();
		}
		function wr(e) {
			if (e === `selectionchange` || e === `keyup` || e === `keydown`) return hr(pr);
		}
		function Tr(e, t) {
			if (e === `click`) return hr(t);
		}
		function Er(e, t) {
			if (e === `input` || e === `change`) return hr(t);
		}
		function Dr(e, t) {
			return (e === t && (e !== 0 || 1 / e == 1 / t)) || (e !== e && t !== t);
		}
		var Or = typeof Object.is == `function` ? Object.is : Dr;
		function kr(e, t) {
			if (Or(e, t)) return !0;
			if (typeof e != `object` || !e || typeof t != `object` || !t) return !1;
			var n = Object.keys(e),
				r = Object.keys(t);
			if (n.length !== r.length) return !1;
			for (r = 0; r < n.length; r++) {
				var i = n[r];
				if (!Oe.call(t, i) || !Or(e[i], t[i])) return !1;
			}
			return !0;
		}
		function Ar(e) {
			for (; e && e.firstChild; ) e = e.firstChild;
			return e;
		}
		function jr(e, t) {
			var n = Ar(e);
			e = 0;
			for (var r; n; ) {
				if (n.nodeType === 3) {
					if (((r = e + n.textContent.length), e <= t && r >= t)) return { node: n, offset: t - e };
					e = r;
				}
				a: {
					for (; n; ) {
						if (n.nextSibling) {
							n = n.nextSibling;
							break a;
						}
						n = n.parentNode;
					}
					n = void 0;
				}
				n = Ar(n);
			}
		}
		function Mr(e, t) {
			return e && t
				? e === t
					? !0
					: e && e.nodeType === 3
						? !1
						: t && t.nodeType === 3
							? Mr(e, t.parentNode)
							: `contains` in e
								? e.contains(t)
								: e.compareDocumentPosition
									? !!(e.compareDocumentPosition(t) & 16)
									: !1
				: !1;
		}
		function Nr(e) {
			e =
				e != null && e.ownerDocument != null && e.ownerDocument.defaultView != null
					? e.ownerDocument.defaultView
					: window;
			for (var t = Wt(e.document); t instanceof e.HTMLIFrameElement; ) {
				try {
					var n = typeof t.contentWindow.location.href == `string`;
				} catch {
					n = !1;
				}
				if (n) e = t.contentWindow;
				else break;
				t = Wt(e.document);
			}
			return t;
		}
		function Pr(e) {
			var t = e && e.nodeName && e.nodeName.toLowerCase();
			return (
				t &&
				((t === `input` &&
					(e.type === `text` ||
						e.type === `search` ||
						e.type === `tel` ||
						e.type === `url` ||
						e.type === `password`)) ||
					t === `textarea` ||
					e.contentEditable === `true`)
			);
		}
		var Fr = gn && `documentMode` in document && 11 >= document.documentMode,
			Ir = null,
			Lr = null,
			Rr = null,
			zr = !1;
		function Br(e, t, n) {
			var r = n.window === n ? n.document : n.nodeType === 9 ? n : n.ownerDocument;
			zr ||
				Ir == null ||
				Ir !== Wt(r) ||
				((r = Ir),
				`selectionStart` in r && Pr(r)
					? (r = { start: r.selectionStart, end: r.selectionEnd })
					: ((r = ((r.ownerDocument && r.ownerDocument.defaultView) || window).getSelection()),
						(r = {
							anchorNode: r.anchorNode,
							anchorOffset: r.anchorOffset,
							focusNode: r.focusNode,
							focusOffset: r.focusOffset,
						})),
				(Rr && kr(Rr, r)) ||
					((Rr = r),
					(r = Ed(Lr, `onSelect`)),
					0 < r.length &&
						((t = new On(`onSelect`, `select`, null, t, n)),
						e.push({ event: t, listeners: r }),
						(t.target = Ir))));
		}
		function Vr(e, t) {
			var n = {};
			return (
				(n[e.toLowerCase()] = t.toLowerCase()),
				(n[`Webkit` + e] = `webkit` + t),
				(n[`Moz` + e] = `moz` + t),
				n
			);
		}
		var Hr = {
				animationend: Vr(`Animation`, `AnimationEnd`),
				animationiteration: Vr(`Animation`, `AnimationIteration`),
				animationstart: Vr(`Animation`, `AnimationStart`),
				transitionrun: Vr(`Transition`, `TransitionRun`),
				transitionstart: Vr(`Transition`, `TransitionStart`),
				transitioncancel: Vr(`Transition`, `TransitionCancel`),
				transitionend: Vr(`Transition`, `TransitionEnd`),
			},
			Ur = {},
			Wr = {};
		gn &&
			((Wr = document.createElement(`div`).style),
			`AnimationEvent` in window ||
				(delete Hr.animationend.animation,
				delete Hr.animationiteration.animation,
				delete Hr.animationstart.animation),
			`TransitionEvent` in window || delete Hr.transitionend.transition);
		function Gr(e) {
			if (Ur[e]) return Ur[e];
			if (!Hr[e]) return e;
			var t = Hr[e],
				n;
			for (n in t) if (t.hasOwnProperty(n) && n in Wr) return (Ur[e] = t[n]);
			return e;
		}
		var Kr = Gr(`animationend`),
			qr = Gr(`animationiteration`),
			Jr = Gr(`animationstart`),
			Yr = Gr(`transitionrun`),
			Xr = Gr(`transitionstart`),
			Zr = Gr(`transitioncancel`),
			Qr = Gr(`transitionend`),
			$r = new Map(),
			ei =
				`abort auxClick beforeToggle cancel canPlay canPlayThrough click close contextMenu copy cut drag dragEnd dragEnter dragExit dragLeave dragOver dragStart drop durationChange emptied encrypted ended error gotPointerCapture input invalid keyDown keyPress keyUp load loadedData loadedMetadata loadStart lostPointerCapture mouseDown mouseMove mouseOut mouseOver mouseUp paste pause play playing pointerCancel pointerDown pointerMove pointerOut pointerOver pointerUp progress rateChange reset resize seeked seeking stalled submit suspend timeUpdate touchCancel touchEnd touchStart volumeChange scroll toggle touchMove waiting wheel`.split(
					` `,
				);
		ei.push(`scrollEnd`);
		function ti(e, t) {
			$r.set(e, t), At(t, [e]);
		}
		var ni =
				typeof reportError == `function`
					? reportError
					: function (e) {
							if (typeof window == `object` && typeof window.ErrorEvent == `function`) {
								var t = new window.ErrorEvent(`error`, {
									bubbles: !0,
									cancelable: !0,
									message:
										typeof e == `object` && e && typeof e.message == `string`
											? String(e.message)
											: String(e),
									error: e,
								});
								if (!window.dispatchEvent(t)) return;
							} else if (typeof process == `object` && typeof process.emit == `function`) {
								process.emit(`uncaughtException`, e);
								return;
							}
							console.error(e);
						},
			ri = [],
			ii = 0,
			ai = 0;
		function oi() {
			for (var e = ii, t = (ai = ii = 0); t < e; ) {
				var n = ri[t];
				ri[t++] = null;
				var r = ri[t];
				ri[t++] = null;
				var i = ri[t];
				ri[t++] = null;
				var a = ri[t];
				if (((ri[t++] = null), r !== null && i !== null)) {
					var o = r.pending;
					o === null ? (i.next = i) : ((i.next = o.next), (o.next = i)), (r.pending = i);
				}
				a !== 0 && ui(n, i, a);
			}
		}
		function si(e, t, n, r) {
			(ri[ii++] = e),
				(ri[ii++] = t),
				(ri[ii++] = n),
				(ri[ii++] = r),
				(ai |= r),
				(e.lanes |= r),
				(e = e.alternate),
				e !== null && (e.lanes |= r);
		}
		function ci(e, t, n, r) {
			return si(e, t, n, r), di(e);
		}
		function li(e, t) {
			return si(e, null, null, t), di(e);
		}
		function ui(e, t, n) {
			e.lanes |= n;
			var r = e.alternate;
			r !== null && (r.lanes |= n);
			for (var i = !1, a = e.return; a !== null; )
				(a.childLanes |= n),
					(r = a.alternate),
					r !== null && (r.childLanes |= n),
					a.tag === 22 && ((e = a.stateNode), e === null || e._visibility & 1 || (i = !0)),
					(e = a),
					(a = a.return);
			return e.tag === 3
				? ((a = e.stateNode),
					i &&
						t !== null &&
						((i = 31 - Ge(n)),
						(e = a.hiddenUpdates),
						(r = e[i]),
						r === null ? (e[i] = [t]) : r.push(t),
						(t.lane = n | 536870912)),
					a)
				: null;
		}
		function di(e) {
			if (50 < du) throw ((du = 0), (fu = null), Error(s(185)));
			for (var t = e.return; t !== null; ) (e = t), (t = e.return);
			return e.tag === 3 ? e.stateNode : null;
		}
		var fi = {};
		function pi(e, t, n, r) {
			(this.tag = e),
				(this.key = n),
				(this.sibling =
					this.child =
					this.return =
					this.stateNode =
					this.type =
					this.elementType =
						null),
				(this.index = 0),
				(this.refCleanup = this.ref = null),
				(this.pendingProps = t),
				(this.dependencies = this.memoizedState = this.updateQueue = this.memoizedProps = null),
				(this.mode = r),
				(this.subtreeFlags = this.flags = 0),
				(this.deletions = null),
				(this.childLanes = this.lanes = 0),
				(this.alternate = null);
		}
		function mi(e, t, n, r) {
			return new pi(e, t, n, r);
		}
		function hi(e) {
			return (e = e.prototype), !(!e || !e.isReactComponent);
		}
		function gi(e, t) {
			var n = e.alternate;
			return (
				n === null
					? ((n = mi(e.tag, t, e.key, e.mode)),
						(n.elementType = e.elementType),
						(n.type = e.type),
						(n.stateNode = e.stateNode),
						(n.alternate = e),
						(e.alternate = n))
					: ((n.pendingProps = t),
						(n.type = e.type),
						(n.flags = 0),
						(n.subtreeFlags = 0),
						(n.deletions = null)),
				(n.flags = e.flags & 65011712),
				(n.childLanes = e.childLanes),
				(n.lanes = e.lanes),
				(n.child = e.child),
				(n.memoizedProps = e.memoizedProps),
				(n.memoizedState = e.memoizedState),
				(n.updateQueue = e.updateQueue),
				(t = e.dependencies),
				(n.dependencies = t === null ? null : { lanes: t.lanes, firstContext: t.firstContext }),
				(n.sibling = e.sibling),
				(n.index = e.index),
				(n.ref = e.ref),
				(n.refCleanup = e.refCleanup),
				n
			);
		}
		function _i(e, t) {
			e.flags &= 65011714;
			var n = e.alternate;
			return (
				n === null
					? ((e.childLanes = 0),
						(e.lanes = t),
						(e.child = null),
						(e.subtreeFlags = 0),
						(e.memoizedProps = null),
						(e.memoizedState = null),
						(e.updateQueue = null),
						(e.dependencies = null),
						(e.stateNode = null))
					: ((e.childLanes = n.childLanes),
						(e.lanes = n.lanes),
						(e.child = n.child),
						(e.subtreeFlags = 0),
						(e.deletions = null),
						(e.memoizedProps = n.memoizedProps),
						(e.memoizedState = n.memoizedState),
						(e.updateQueue = n.updateQueue),
						(e.type = n.type),
						(t = n.dependencies),
						(e.dependencies =
							t === null ? null : { lanes: t.lanes, firstContext: t.firstContext })),
				e
			);
		}
		function vi(e, t, n, r, i, a) {
			var o = 0;
			if (((r = e), typeof e == `function`)) hi(e) && (o = 1);
			else if (typeof e == `string`)
				o = Uf(e, n, me.current) ? 26 : e === `html` || e === `head` || e === `body` ? 27 : 5;
			else
				a: switch (e) {
					case ie:
						return (e = mi(31, n, t, i)), (e.elementType = ie), (e.lanes = a), e;
					case y:
						return yi(n.children, i, a, t);
					case b:
						(o = 8), (i |= 24);
						break;
					case x:
						return (e = mi(12, n, t, i | 2)), (e.elementType = x), (e.lanes = a), e;
					case te:
						return (e = mi(13, n, t, i)), (e.elementType = te), (e.lanes = a), e;
					case ne:
						return (e = mi(19, n, t, i)), (e.elementType = ne), (e.lanes = a), e;
					default:
						if (typeof e == `object` && e)
							switch (e.$$typeof) {
								case S:
									o = 10;
									break a;
								case ee:
									o = 9;
									break a;
								case C:
									o = 11;
									break a;
								case re:
									o = 14;
									break a;
								case w:
									(o = 16), (r = null);
									break a;
							}
						(o = 29), (n = Error(s(130, e === null ? `null` : typeof e, ``))), (r = null);
				}
			return (t = mi(o, n, t, i)), (t.elementType = e), (t.type = r), (t.lanes = a), t;
		}
		function yi(e, t, n, r) {
			return (e = mi(7, e, r, t)), (e.lanes = n), e;
		}
		function bi(e, t, n) {
			return (e = mi(6, e, null, t)), (e.lanes = n), e;
		}
		function xi(e) {
			var t = mi(18, null, null, 0);
			return (t.stateNode = e), t;
		}
		function Si(e, t, n) {
			return (
				(t = mi(4, e.children === null ? [] : e.children, e.key, t)),
				(t.lanes = n),
				(t.stateNode = {
					containerInfo: e.containerInfo,
					pendingChildren: null,
					implementation: e.implementation,
				}),
				t
			);
		}
		var Ci = new WeakMap();
		function wi(e, t) {
			if (typeof e == `object` && e) {
				var n = Ci.get(e);
				return n === void 0 ? ((t = { value: e, source: t, stack: A(t) }), Ci.set(e, t), t) : n;
			}
			return { value: e, source: t, stack: A(t) };
		}
		var Ti = [],
			Ei = 0,
			Di = null,
			Oi = 0,
			ki = [],
			Ai = 0,
			ji = null,
			Mi = 1,
			Ni = ``;
		function Pi(e, t) {
			(Ti[Ei++] = Oi), (Ti[Ei++] = Di), (Di = e), (Oi = t);
		}
		function Fi(e, t, n) {
			(ki[Ai++] = Mi), (ki[Ai++] = Ni), (ki[Ai++] = ji), (ji = e);
			var r = Mi;
			e = Ni;
			var i = 32 - Ge(r) - 1;
			(r &= ~(1 << i)), (n += 1);
			var a = 32 - Ge(t) + i;
			if (30 < a) {
				var o = i - (i % 5);
				(a = (r & ((1 << o) - 1)).toString(32)),
					(r >>= o),
					(i -= o),
					(Mi = (1 << (32 - Ge(t) + i)) | (n << i) | r),
					(Ni = a + e);
			} else (Mi = (1 << a) | (n << i) | r), (Ni = e);
		}
		function Ii(e) {
			e.return !== null && (Pi(e, 1), Fi(e, 1, 0));
		}
		function Li(e) {
			for (; e === Di; ) (Di = Ti[--Ei]), (Ti[Ei] = null), (Oi = Ti[--Ei]), (Ti[Ei] = null);
			for (; e === ji; )
				(ji = ki[--Ai]),
					(ki[Ai] = null),
					(Ni = ki[--Ai]),
					(ki[Ai] = null),
					(Mi = ki[--Ai]),
					(ki[Ai] = null);
		}
		function Ri(e, t) {
			(ki[Ai++] = Mi), (ki[Ai++] = Ni), (ki[Ai++] = ji), (Mi = t.id), (Ni = t.overflow), (ji = e);
		}
		var zi = null,
			N = null,
			P = !1,
			Bi = null,
			Vi = !1,
			Hi = Error(s(519));
		function Ui(e) {
			throw (
				(Yi(
					wi(
						Error(
							s(
								418,
								1 < arguments.length && arguments[1] !== void 0 && arguments[1] ? `text` : `HTML`,
								``,
							),
						),
						e,
					),
				),
				Hi)
			);
		}
		function Wi(e) {
			var t = e.stateNode,
				n = e.type,
				r = e.memoizedProps;
			switch (((t[mt] = e), (t[ht] = r), n)) {
				case `dialog`:
					Q(`cancel`, t), Q(`close`, t);
					break;
				case `iframe`:
				case `object`:
				case `embed`:
					Q(`load`, t);
					break;
				case `video`:
				case `audio`:
					for (n = 0; n < _d.length; n++) Q(_d[n], t);
					break;
				case `source`:
					Q(`error`, t);
					break;
				case `img`:
				case `image`:
				case `link`:
					Q(`error`, t), Q(`load`, t);
					break;
				case `details`:
					Q(`toggle`, t);
					break;
				case `input`:
					Q(`invalid`, t),
						Jt(t, r.value, r.defaultValue, r.checked, r.defaultChecked, r.type, r.name, !0);
					break;
				case `select`:
					Q(`invalid`, t);
					break;
				case `textarea`:
					Q(`invalid`, t), Qt(t, r.value, r.defaultValue, r.children);
			}
			(n = r.children),
				(typeof n != `string` && typeof n != `number` && typeof n != `bigint`) ||
				t.textContent === `` + n ||
				!0 === r.suppressHydrationWarning ||
				Md(t.textContent, n)
					? (r.popover != null && (Q(`beforetoggle`, t), Q(`toggle`, t)),
						r.onScroll != null && Q(`scroll`, t),
						r.onScrollEnd != null && Q(`scrollend`, t),
						r.onClick != null && (t.onclick = sn),
						(t = !0))
					: (t = !1),
				t || Ui(e, !0);
		}
		function Gi(e) {
			for (zi = e.return; zi; )
				switch (zi.tag) {
					case 5:
					case 31:
					case 13:
						Vi = !1;
						return;
					case 27:
					case 3:
						Vi = !0;
						return;
					default:
						zi = zi.return;
				}
		}
		function Ki(e) {
			if (e !== zi) return !1;
			if (!P) return Gi(e), (P = !0), !1;
			var t = e.tag,
				n;
			if (
				((n = t !== 3 && t !== 27) &&
					((n = t === 5) &&
						((n = e.type), (n = !(n !== `form` && n !== `button`) || Ud(e.type, e.memoizedProps))),
					(n = !n)),
				n && N && Ui(e),
				Gi(e),
				t === 13)
			) {
				if (((e = e.memoizedState), (e = e === null ? null : e.dehydrated), !e))
					throw Error(s(317));
				N = uf(e);
			} else if (t === 31) {
				if (((e = e.memoizedState), (e = e === null ? null : e.dehydrated), !e))
					throw Error(s(317));
				N = uf(e);
			} else
				t === 27
					? ((t = N), Zd(e.type) ? ((e = lf), (lf = null), (N = e)) : (N = t))
					: (N = zi ? cf(e.stateNode.nextSibling) : null);
			return !0;
		}
		function qi() {
			(N = zi = null), (P = !1);
		}
		function Ji() {
			var e = Bi;
			return e !== null && (Zl === null ? (Zl = e) : Zl.push.apply(Zl, e), (Bi = null)), e;
		}
		function Yi(e) {
			Bi === null ? (Bi = [e]) : Bi.push(e);
		}
		var Xi = D(null),
			Zi = null,
			Qi = null;
		function $i(e, t, n) {
			k(Xi, t._currentValue), (t._currentValue = n);
		}
		function ea(e) {
			(e._currentValue = Xi.current), O(Xi);
		}
		function ta(e, t, n) {
			for (; e !== null; ) {
				var r = e.alternate;
				if (
					((e.childLanes & t) === t
						? r !== null && (r.childLanes & t) !== t && (r.childLanes |= t)
						: ((e.childLanes |= t), r !== null && (r.childLanes |= t)),
					e === n)
				)
					break;
				e = e.return;
			}
		}
		function na(e, t, n, r) {
			var i = e.child;
			for (i !== null && (i.return = e); i !== null; ) {
				var a = i.dependencies;
				if (a !== null) {
					var o = i.child;
					a = a.firstContext;
					a: for (; a !== null; ) {
						var c = a;
						a = i;
						for (var l = 0; l < t.length; l++)
							if (c.context === t[l]) {
								(a.lanes |= n),
									(c = a.alternate),
									c !== null && (c.lanes |= n),
									ta(a.return, n, e),
									r || (o = null);
								break a;
							}
						a = c.next;
					}
				} else if (i.tag === 18) {
					if (((o = i.return), o === null)) throw Error(s(341));
					(o.lanes |= n), (a = o.alternate), a !== null && (a.lanes |= n), ta(o, n, e), (o = null);
				} else o = i.child;
				if (o !== null) o.return = i;
				else
					for (o = i; o !== null; ) {
						if (o === e) {
							o = null;
							break;
						}
						if (((i = o.sibling), i !== null)) {
							(i.return = o.return), (o = i);
							break;
						}
						o = o.return;
					}
				i = o;
			}
		}
		function ra(e, t, n, r) {
			e = null;
			for (var i = t, a = !1; i !== null; ) {
				if (!a) {
					if (i.flags & 524288) a = !0;
					else if (i.flags & 262144) break;
				}
				if (i.tag === 10) {
					var o = i.alternate;
					if (o === null) throw Error(s(387));
					if (((o = o.memoizedProps), o !== null)) {
						var c = i.type;
						Or(i.pendingProps.value, o.value) || (e === null ? (e = [c]) : e.push(c));
					}
				} else if (i === _e.current) {
					if (((o = i.alternate), o === null)) throw Error(s(387));
					o.memoizedState.memoizedState !== i.memoizedState.memoizedState &&
						(e === null ? (e = [Qf]) : e.push(Qf));
				}
				i = i.return;
			}
			e !== null && na(t, e, n, r), (t.flags |= 262144);
		}
		function ia(e) {
			for (e = e.firstContext; e !== null; ) {
				if (!Or(e.context._currentValue, e.memoizedValue)) return !0;
				e = e.next;
			}
			return !1;
		}
		function aa(e) {
			(Zi = e), (Qi = null), (e = e.dependencies), e !== null && (e.firstContext = null);
		}
		function oa(e) {
			return ca(Zi, e);
		}
		function sa(e, t) {
			return Zi === null && aa(e), ca(e, t);
		}
		function ca(e, t) {
			var n = t._currentValue;
			if (((t = { context: t, memoizedValue: n, next: null }), Qi === null)) {
				if (e === null) throw Error(s(308));
				(Qi = t), (e.dependencies = { lanes: 0, firstContext: t }), (e.flags |= 524288);
			} else Qi = Qi.next = t;
			return n;
		}
		var la =
				typeof AbortController < `u`
					? AbortController
					: function () {
							var e = [],
								t = (this.signal = {
									aborted: !1,
									addEventListener: function (t, n) {
										e.push(n);
									},
								});
							this.abort = function () {
								(t.aborted = !0),
									e.forEach(function (e) {
										return e();
									});
							};
						},
			ua = t.unstable_scheduleCallback,
			da = t.unstable_NormalPriority,
			F = {
				$$typeof: S,
				Consumer: null,
				Provider: null,
				_currentValue: null,
				_currentValue2: null,
				_threadCount: 0,
			};
		function fa() {
			return { controller: new la(), data: new Map(), refCount: 0 };
		}
		function pa(e) {
			e.refCount--,
				e.refCount === 0 &&
					ua(da, function () {
						e.controller.abort();
					});
		}
		var ma = null,
			ha = 0,
			ga = 0,
			_a = null;
		function va(e, t) {
			if (ma === null) {
				var n = (ma = []);
				(ha = 0),
					(ga = dd()),
					(_a = {
						status: `pending`,
						value: void 0,
						then: function (e) {
							n.push(e);
						},
					});
			}
			return ha++, t.then(ya, ya), t;
		}
		function ya() {
			if (--ha === 0 && ma !== null) {
				_a !== null && (_a.status = `fulfilled`);
				var e = ma;
				(ma = null), (ga = 0), (_a = null);
				for (var t = 0; t < e.length; t++) (0, e[t])();
			}
		}
		function ba(e, t) {
			var n = [],
				r = {
					status: `pending`,
					value: null,
					reason: null,
					then: function (e) {
						n.push(e);
					},
				};
			return (
				e.then(
					function () {
						(r.status = `fulfilled`), (r.value = t);
						for (var e = 0; e < n.length; e++) (0, n[e])(t);
					},
					function (e) {
						for (r.status = `rejected`, r.reason = e, e = 0; e < n.length; e++) (0, n[e])(void 0);
					},
				),
				r
			);
		}
		var xa = T.S;
		T.S = function (e, t) {
			(eu = Ne()),
				typeof t == `object` && t && typeof t.then == `function` && va(e, t),
				xa !== null && xa(e, t);
		};
		var Sa = D(null);
		function Ca() {
			var e = Sa.current;
			return e === null ? K.pooledCache : e;
		}
		function wa(e, t) {
			t === null ? k(Sa, Sa.current) : k(Sa, t.pool);
		}
		function Ta() {
			var e = Ca();
			return e === null ? null : { parent: F._currentValue, pool: e };
		}
		var Ea = Error(s(460)),
			Da = Error(s(474)),
			Oa = Error(s(542)),
			ka = { then: function () {} };
		function Aa(e) {
			return (e = e.status), e === `fulfilled` || e === `rejected`;
		}
		function ja(e, t, n) {
			switch (
				((n = e[n]), n === void 0 ? e.push(t) : n !== t && (t.then(sn, sn), (t = n)), t.status)
			) {
				case `fulfilled`:
					return t.value;
				case `rejected`:
					throw ((e = t.reason), Fa(e), e);
				default:
					if (typeof t.status == `string`) t.then(sn, sn);
					else {
						if (((e = K), e !== null && 100 < e.shellSuspendCounter)) throw Error(s(482));
						(e = t),
							(e.status = `pending`),
							e.then(
								function (e) {
									if (t.status === `pending`) {
										var n = t;
										(n.status = `fulfilled`), (n.value = e);
									}
								},
								function (e) {
									if (t.status === `pending`) {
										var n = t;
										(n.status = `rejected`), (n.reason = e);
									}
								},
							);
					}
					switch (t.status) {
						case `fulfilled`:
							return t.value;
						case `rejected`:
							throw ((e = t.reason), Fa(e), e);
					}
					throw ((Na = t), Ea);
			}
		}
		function Ma(e) {
			try {
				var t = e._init;
				return t(e._payload);
			} catch (e) {
				throw typeof e == `object` && e && typeof e.then == `function` ? ((Na = e), Ea) : e;
			}
		}
		var Na = null;
		function Pa() {
			if (Na === null) throw Error(s(459));
			var e = Na;
			return (Na = null), e;
		}
		function Fa(e) {
			if (e === Ea || e === Oa) throw Error(s(483));
		}
		var Ia = null,
			La = 0;
		function Ra(e) {
			var t = La;
			return (La += 1), Ia === null && (Ia = []), ja(Ia, e, t);
		}
		function za(e, t) {
			(t = t.props.ref), (e.ref = t === void 0 ? null : t);
		}
		function Ba(e, t) {
			throw t.$$typeof === g
				? Error(s(525))
				: ((e = Object.prototype.toString.call(t)),
					Error(
						s(
							31,
							e === `[object Object]` ? `object with keys {` + Object.keys(t).join(`, `) + `}` : e,
						),
					));
		}
		function Va(e) {
			function t(t, n) {
				if (e) {
					var r = t.deletions;
					r === null ? ((t.deletions = [n]), (t.flags |= 16)) : r.push(n);
				}
			}
			function n(n, r) {
				if (!e) return null;
				for (; r !== null; ) t(n, r), (r = r.sibling);
				return null;
			}
			function r(e) {
				for (var t = new Map(); e !== null; )
					e.key === null ? t.set(e.index, e) : t.set(e.key, e), (e = e.sibling);
				return t;
			}
			function i(e, t) {
				return (e = gi(e, t)), (e.index = 0), (e.sibling = null), e;
			}
			function a(t, n, r) {
				return (
					(t.index = r),
					e
						? ((r = t.alternate),
							r === null
								? ((t.flags |= 67108866), n)
								: ((r = r.index), r < n ? ((t.flags |= 67108866), n) : r))
						: ((t.flags |= 1048576), n)
				);
			}
			function o(t) {
				return e && t.alternate === null && (t.flags |= 67108866), t;
			}
			function c(e, t, n, r) {
				return t === null || t.tag !== 6
					? ((t = bi(n, e.mode, r)), (t.return = e), t)
					: ((t = i(t, n)), (t.return = e), t);
			}
			function l(e, t, n, r) {
				var a = n.type;
				return a === y
					? d(e, t, n.props.children, r, n.key)
					: t !== null &&
							(t.elementType === a ||
								(typeof a == `object` && a && a.$$typeof === w && Ma(a) === t.type))
						? ((t = i(t, n.props)), za(t, n), (t.return = e), t)
						: ((t = vi(n.type, n.key, n.props, null, e.mode, r)), za(t, n), (t.return = e), t);
			}
			function u(e, t, n, r) {
				return t === null ||
					t.tag !== 4 ||
					t.stateNode.containerInfo !== n.containerInfo ||
					t.stateNode.implementation !== n.implementation
					? ((t = Si(n, e.mode, r)), (t.return = e), t)
					: ((t = i(t, n.children || [])), (t.return = e), t);
			}
			function d(e, t, n, r, a) {
				return t === null || t.tag !== 7
					? ((t = yi(n, e.mode, r, a)), (t.return = e), t)
					: ((t = i(t, n)), (t.return = e), t);
			}
			function f(e, t, n) {
				if ((typeof t == `string` && t !== ``) || typeof t == `number` || typeof t == `bigint`)
					return (t = bi(`` + t, e.mode, n)), (t.return = e), t;
				if (typeof t == `object` && t) {
					switch (t.$$typeof) {
						case _:
							return (n = vi(t.type, t.key, t.props, null, e.mode, n)), za(n, t), (n.return = e), n;
						case v:
							return (t = Si(t, e.mode, n)), (t.return = e), t;
						case w:
							return (t = Ma(t)), f(e, t, n);
					}
					if (ue(t) || se(t)) return (t = yi(t, e.mode, n, null)), (t.return = e), t;
					if (typeof t.then == `function`) return f(e, Ra(t), n);
					if (t.$$typeof === S) return f(e, sa(e, t), n);
					Ba(e, t);
				}
				return null;
			}
			function p(e, t, n, r) {
				var i = t === null ? null : t.key;
				if ((typeof n == `string` && n !== ``) || typeof n == `number` || typeof n == `bigint`)
					return i === null ? c(e, t, `` + n, r) : null;
				if (typeof n == `object` && n) {
					switch (n.$$typeof) {
						case _:
							return n.key === i ? l(e, t, n, r) : null;
						case v:
							return n.key === i ? u(e, t, n, r) : null;
						case w:
							return (n = Ma(n)), p(e, t, n, r);
					}
					if (ue(n) || se(n)) return i === null ? d(e, t, n, r, null) : null;
					if (typeof n.then == `function`) return p(e, t, Ra(n), r);
					if (n.$$typeof === S) return p(e, t, sa(e, n), r);
					Ba(e, n);
				}
				return null;
			}
			function m(e, t, n, r, i) {
				if ((typeof r == `string` && r !== ``) || typeof r == `number` || typeof r == `bigint`)
					return (e = e.get(n) || null), c(t, e, `` + r, i);
				if (typeof r == `object` && r) {
					switch (r.$$typeof) {
						case _:
							return (e = e.get(r.key === null ? n : r.key) || null), l(t, e, r, i);
						case v:
							return (e = e.get(r.key === null ? n : r.key) || null), u(t, e, r, i);
						case w:
							return (r = Ma(r)), m(e, t, n, r, i);
					}
					if (ue(r) || se(r)) return (e = e.get(n) || null), d(t, e, r, i, null);
					if (typeof r.then == `function`) return m(e, t, n, Ra(r), i);
					if (r.$$typeof === S) return m(e, t, n, sa(t, r), i);
					Ba(t, r);
				}
				return null;
			}
			function h(i, o, s, c) {
				for (
					var l = null, u = null, d = o, h = (o = 0), g = null;
					d !== null && h < s.length;
					h++
				) {
					d.index > h ? ((g = d), (d = null)) : (g = d.sibling);
					var _ = p(i, d, s[h], c);
					if (_ === null) {
						d === null && (d = g);
						break;
					}
					e && d && _.alternate === null && t(i, d),
						(o = a(_, o, h)),
						u === null ? (l = _) : (u.sibling = _),
						(u = _),
						(d = g);
				}
				if (h === s.length) return n(i, d), P && Pi(i, h), l;
				if (d === null) {
					for (; h < s.length; h++)
						(d = f(i, s[h], c)),
							d !== null && ((o = a(d, o, h)), u === null ? (l = d) : (u.sibling = d), (u = d));
					return P && Pi(i, h), l;
				}
				for (d = r(d); h < s.length; h++)
					(g = m(d, i, h, s[h], c)),
						g !== null &&
							(e && g.alternate !== null && d.delete(g.key === null ? h : g.key),
							(o = a(g, o, h)),
							u === null ? (l = g) : (u.sibling = g),
							(u = g));
				return (
					e &&
						d.forEach(function (e) {
							return t(i, e);
						}),
					P && Pi(i, h),
					l
				);
			}
			function g(i, o, c, l) {
				if (c == null) throw Error(s(151));
				for (
					var u = null, d = null, h = o, g = (o = 0), _ = null, v = c.next();
					h !== null && !v.done;
					g++, v = c.next()
				) {
					h.index > g ? ((_ = h), (h = null)) : (_ = h.sibling);
					var y = p(i, h, v.value, l);
					if (y === null) {
						h === null && (h = _);
						break;
					}
					e && h && y.alternate === null && t(i, h),
						(o = a(y, o, g)),
						d === null ? (u = y) : (d.sibling = y),
						(d = y),
						(h = _);
				}
				if (v.done) return n(i, h), P && Pi(i, g), u;
				if (h === null) {
					for (; !v.done; g++, v = c.next())
						(v = f(i, v.value, l)),
							v !== null && ((o = a(v, o, g)), d === null ? (u = v) : (d.sibling = v), (d = v));
					return P && Pi(i, g), u;
				}
				for (h = r(h); !v.done; g++, v = c.next())
					(v = m(h, i, g, v.value, l)),
						v !== null &&
							(e && v.alternate !== null && h.delete(v.key === null ? g : v.key),
							(o = a(v, o, g)),
							d === null ? (u = v) : (d.sibling = v),
							(d = v));
				return (
					e &&
						h.forEach(function (e) {
							return t(i, e);
						}),
					P && Pi(i, g),
					u
				);
			}
			function b(e, r, a, c) {
				if (
					(typeof a == `object` && a && a.type === y && a.key === null && (a = a.props.children),
					typeof a == `object` && a)
				) {
					switch (a.$$typeof) {
						case _:
							a: {
								for (var l = a.key; r !== null; ) {
									if (r.key === l) {
										if (((l = a.type), l === y)) {
											if (r.tag === 7) {
												n(e, r.sibling), (c = i(r, a.props.children)), (c.return = e), (e = c);
												break a;
											}
										} else if (
											r.elementType === l ||
											(typeof l == `object` && l && l.$$typeof === w && Ma(l) === r.type)
										) {
											n(e, r.sibling), (c = i(r, a.props)), za(c, a), (c.return = e), (e = c);
											break a;
										}
										n(e, r);
										break;
									} else t(e, r);
									r = r.sibling;
								}
								a.type === y
									? ((c = yi(a.props.children, e.mode, c, a.key)), (c.return = e), (e = c))
									: ((c = vi(a.type, a.key, a.props, null, e.mode, c)),
										za(c, a),
										(c.return = e),
										(e = c));
							}
							return o(e);
						case v:
							a: {
								for (l = a.key; r !== null; ) {
									if (r.key === l)
										if (
											r.tag === 4 &&
											r.stateNode.containerInfo === a.containerInfo &&
											r.stateNode.implementation === a.implementation
										) {
											n(e, r.sibling), (c = i(r, a.children || [])), (c.return = e), (e = c);
											break a;
										} else {
											n(e, r);
											break;
										}
									else t(e, r);
									r = r.sibling;
								}
								(c = Si(a, e.mode, c)), (c.return = e), (e = c);
							}
							return o(e);
						case w:
							return (a = Ma(a)), b(e, r, a, c);
					}
					if (ue(a)) return h(e, r, a, c);
					if (se(a)) {
						if (((l = se(a)), typeof l != `function`)) throw Error(s(150));
						return (a = l.call(a)), g(e, r, a, c);
					}
					if (typeof a.then == `function`) return b(e, r, Ra(a), c);
					if (a.$$typeof === S) return b(e, r, sa(e, a), c);
					Ba(e, a);
				}
				return (typeof a == `string` && a !== ``) || typeof a == `number` || typeof a == `bigint`
					? ((a = `` + a),
						r !== null && r.tag === 6
							? (n(e, r.sibling), (c = i(r, a)), (c.return = e), (e = c))
							: (n(e, r), (c = bi(a, e.mode, c)), (c.return = e), (e = c)),
						o(e))
					: n(e, r);
			}
			return function (e, t, n, r) {
				try {
					La = 0;
					var i = b(e, t, n, r);
					return (Ia = null), i;
				} catch (t) {
					if (t === Ea || t === Oa) throw t;
					var a = mi(29, t, null, e.mode);
					return (a.lanes = r), (a.return = e), a;
				}
			};
		}
		var Ha = Va(!0),
			Ua = Va(!1),
			Wa = !1;
		function Ga(e) {
			e.updateQueue = {
				baseState: e.memoizedState,
				firstBaseUpdate: null,
				lastBaseUpdate: null,
				shared: { pending: null, lanes: 0, hiddenCallbacks: null },
				callbacks: null,
			};
		}
		function Ka(e, t) {
			(e = e.updateQueue),
				t.updateQueue === e &&
					(t.updateQueue = {
						baseState: e.baseState,
						firstBaseUpdate: e.firstBaseUpdate,
						lastBaseUpdate: e.lastBaseUpdate,
						shared: e.shared,
						callbacks: null,
					});
		}
		function qa(e) {
			return { lane: e, tag: 0, payload: null, callback: null, next: null };
		}
		function Ja(e, t, n) {
			var r = e.updateQueue;
			if (r === null) return null;
			if (((r = r.shared), G & 2)) {
				var i = r.pending;
				return (
					i === null ? (t.next = t) : ((t.next = i.next), (i.next = t)),
					(r.pending = t),
					(t = di(e)),
					ui(e, null, n),
					t
				);
			}
			return si(e, r, t, n), di(e);
		}
		function Ya(e, t, n) {
			if (((t = t.updateQueue), t !== null && ((t = t.shared), n & 4194048))) {
				var r = t.lanes;
				(r &= e.pendingLanes), (n |= r), (t.lanes = n), st(e, n);
			}
		}
		function Xa(e, t) {
			var n = e.updateQueue,
				r = e.alternate;
			if (r !== null && ((r = r.updateQueue), n === r)) {
				var i = null,
					a = null;
				if (((n = n.firstBaseUpdate), n !== null)) {
					do {
						var o = { lane: n.lane, tag: n.tag, payload: n.payload, callback: null, next: null };
						a === null ? (i = a = o) : (a = a.next = o), (n = n.next);
					} while (n !== null);
					a === null ? (i = a = t) : (a = a.next = t);
				} else i = a = t;
				(n = {
					baseState: r.baseState,
					firstBaseUpdate: i,
					lastBaseUpdate: a,
					shared: r.shared,
					callbacks: r.callbacks,
				}),
					(e.updateQueue = n);
				return;
			}
			(e = n.lastBaseUpdate),
				e === null ? (n.firstBaseUpdate = t) : (e.next = t),
				(n.lastBaseUpdate = t);
		}
		var Za = !1;
		function Qa() {
			if (Za) {
				var e = _a;
				if (e !== null) throw e;
			}
		}
		function $a(e, t, n, r) {
			Za = !1;
			var i = e.updateQueue;
			Wa = !1;
			var a = i.firstBaseUpdate,
				o = i.lastBaseUpdate,
				s = i.shared.pending;
			if (s !== null) {
				i.shared.pending = null;
				var c = s,
					l = c.next;
				(c.next = null), o === null ? (a = l) : (o.next = l), (o = c);
				var u = e.alternate;
				u !== null &&
					((u = u.updateQueue),
					(s = u.lastBaseUpdate),
					s !== o && (s === null ? (u.firstBaseUpdate = l) : (s.next = l), (u.lastBaseUpdate = c)));
			}
			if (a !== null) {
				var d = i.baseState;
				(o = 0), (u = l = c = null), (s = a);
				do {
					var f = s.lane & -536870913,
						p = f !== s.lane;
					if (p ? (J & f) === f : (r & f) === f) {
						f !== 0 && f === ga && (Za = !0),
							u !== null &&
								(u = u.next =
									{ lane: 0, tag: s.tag, payload: s.payload, callback: null, next: null });
						a: {
							var m = e,
								g = s;
							f = t;
							var _ = n;
							switch (g.tag) {
								case 1:
									if (((m = g.payload), typeof m == `function`)) {
										d = m.call(_, d, f);
										break a;
									}
									d = m;
									break a;
								case 3:
									m.flags = (m.flags & -65537) | 128;
								case 0:
									if (
										((m = g.payload), (f = typeof m == `function` ? m.call(_, d, f) : m), f == null)
									)
										break a;
									d = h({}, d, f);
									break a;
								case 2:
									Wa = !0;
							}
						}
						(f = s.callback),
							f !== null &&
								((e.flags |= 64),
								p && (e.flags |= 8192),
								(p = i.callbacks),
								p === null ? (i.callbacks = [f]) : p.push(f));
					} else
						(p = { lane: f, tag: s.tag, payload: s.payload, callback: s.callback, next: null }),
							u === null ? ((l = u = p), (c = d)) : (u = u.next = p),
							(o |= f);
					if (((s = s.next), s === null)) {
						if (((s = i.shared.pending), s === null)) break;
						(p = s),
							(s = p.next),
							(p.next = null),
							(i.lastBaseUpdate = p),
							(i.shared.pending = null);
					}
				} while (1);
				u === null && (c = d),
					(i.baseState = c),
					(i.firstBaseUpdate = l),
					(i.lastBaseUpdate = u),
					a === null && (i.shared.lanes = 0),
					(Gl |= o),
					(e.lanes = o),
					(e.memoizedState = d);
			}
		}
		function eo(e, t) {
			if (typeof e != `function`) throw Error(s(191, e));
			e.call(t);
		}
		function to(e, t) {
			var n = e.callbacks;
			if (n !== null) for (e.callbacks = null, e = 0; e < n.length; e++) eo(n[e], t);
		}
		var no = D(null),
			ro = D(0);
		function io(e, t) {
			(e = Wl), k(ro, e), k(no, t), (Wl = e | t.baseLanes);
		}
		function ao() {
			k(ro, Wl), k(no, no.current);
		}
		function oo() {
			(Wl = ro.current), O(no), O(ro);
		}
		var so = D(null),
			co = null;
		function lo(e) {
			var t = e.alternate;
			k(I, I.current & 1),
				k(so, e),
				co === null && (t === null || no.current !== null || t.memoizedState !== null) && (co = e);
		}
		function uo(e) {
			k(I, I.current), k(so, e), co === null && (co = e);
		}
		function fo(e) {
			e.tag === 22 ? (k(I, I.current), k(so, e), co === null && (co = e)) : po(e);
		}
		function po() {
			k(I, I.current), k(so, so.current);
		}
		function mo(e) {
			O(so), co === e && (co = null), O(I);
		}
		var I = D(0);
		function ho(e) {
			for (var t = e; t !== null; ) {
				if (t.tag === 13) {
					var n = t.memoizedState;
					if (n !== null && ((n = n.dehydrated), n === null || af(n) || of(n))) return t;
				} else if (
					t.tag === 19 &&
					(t.memoizedProps.revealOrder === `forwards` ||
						t.memoizedProps.revealOrder === `backwards` ||
						t.memoizedProps.revealOrder === `unstable_legacy-backwards` ||
						t.memoizedProps.revealOrder === `together`)
				) {
					if (t.flags & 128) return t;
				} else if (t.child !== null) {
					(t.child.return = t), (t = t.child);
					continue;
				}
				if (t === e) break;
				for (; t.sibling === null; ) {
					if (t.return === null || t.return === e) return null;
					t = t.return;
				}
				(t.sibling.return = t.return), (t = t.sibling);
			}
			return null;
		}
		var go = 0,
			L = null,
			R = null,
			z = null,
			_o = !1,
			vo = !1,
			yo = !1,
			bo = 0,
			xo = 0,
			So = null,
			Co = 0;
		function B() {
			throw Error(s(321));
		}
		function wo(e, t) {
			if (t === null) return !1;
			for (var n = 0; n < t.length && n < e.length; n++) if (!Or(e[n], t[n])) return !1;
			return !0;
		}
		function To(e, t, n, r, i, a) {
			return (
				(go = a),
				(L = t),
				(t.memoizedState = null),
				(t.updateQueue = null),
				(t.lanes = 0),
				(T.H = e === null || e.memoizedState === null ? Vs : Hs),
				(yo = !1),
				(a = n(r, i)),
				(yo = !1),
				vo && (a = Do(t, n, r, i)),
				Eo(e),
				a
			);
		}
		function Eo(e) {
			T.H = Bs;
			var t = R !== null && R.next !== null;
			if (((go = 0), (z = R = L = null), (_o = !1), (xo = 0), (So = null), t)) throw Error(s(300));
			e === null || H || ((e = e.dependencies), e !== null && ia(e) && (H = !0));
		}
		function Do(e, t, n, r) {
			L = e;
			var i = 0;
			do {
				if ((vo && (So = null), (xo = 0), (vo = !1), 25 <= i)) throw Error(s(301));
				if (((i += 1), (z = R = null), e.updateQueue != null)) {
					var a = e.updateQueue;
					(a.lastEffect = null),
						(a.events = null),
						(a.stores = null),
						a.memoCache != null && (a.memoCache.index = 0);
				}
				(T.H = Us), (a = t(n, r));
			} while (vo);
			return a;
		}
		function Oo() {
			var e = T.H,
				t = e.useState()[0];
			return (
				(t = typeof t.then == `function` ? Po(t) : t),
				(e = e.useState()[0]),
				(R === null ? null : R.memoizedState) !== e && (L.flags |= 1024),
				t
			);
		}
		function ko() {
			var e = bo !== 0;
			return (bo = 0), e;
		}
		function Ao(e, t, n) {
			(t.updateQueue = e.updateQueue), (t.flags &= -2053), (e.lanes &= ~n);
		}
		function jo(e) {
			if (_o) {
				for (e = e.memoizedState; e !== null; ) {
					var t = e.queue;
					t !== null && (t.pending = null), (e = e.next);
				}
				_o = !1;
			}
			(go = 0), (z = R = L = null), (vo = !1), (xo = bo = 0), (So = null);
		}
		function Mo() {
			var e = { memoizedState: null, baseState: null, baseQueue: null, queue: null, next: null };
			return z === null ? (L.memoizedState = z = e) : (z = z.next = e), z;
		}
		function V() {
			if (R === null) {
				var e = L.alternate;
				e = e === null ? null : e.memoizedState;
			} else e = R.next;
			var t = z === null ? L.memoizedState : z.next;
			if (t !== null) (z = t), (R = e);
			else {
				if (e === null) throw L.alternate === null ? Error(s(467)) : Error(s(310));
				(R = e),
					(e = {
						memoizedState: R.memoizedState,
						baseState: R.baseState,
						baseQueue: R.baseQueue,
						queue: R.queue,
						next: null,
					}),
					z === null ? (L.memoizedState = z = e) : (z = z.next = e);
			}
			return z;
		}
		function No() {
			return { lastEffect: null, events: null, stores: null, memoCache: null };
		}
		function Po(e) {
			var t = xo;
			return (
				(xo += 1),
				So === null && (So = []),
				(e = ja(So, e, t)),
				(t = L),
				(z === null ? t.memoizedState : z.next) === null &&
					((t = t.alternate), (T.H = t === null || t.memoizedState === null ? Vs : Hs)),
				e
			);
		}
		function Fo(e) {
			if (typeof e == `object` && e) {
				if (typeof e.then == `function`) return Po(e);
				if (e.$$typeof === S) return oa(e);
			}
			throw Error(s(438, String(e)));
		}
		function Io(e) {
			var t = null,
				n = L.updateQueue;
			if ((n !== null && (t = n.memoCache), t == null)) {
				var r = L.alternate;
				r !== null &&
					((r = r.updateQueue),
					r !== null &&
						((r = r.memoCache),
						r != null &&
							(t = {
								data: r.data.map(function (e) {
									return e.slice();
								}),
								index: 0,
							})));
			}
			if (
				((t ??= { data: [], index: 0 }),
				n === null && ((n = No()), (L.updateQueue = n)),
				(n.memoCache = t),
				(n = t.data[t.index]),
				n === void 0)
			)
				for (n = t.data[t.index] = Array(e), r = 0; r < e; r++) n[r] = ae;
			return t.index++, n;
		}
		function Lo(e, t) {
			return typeof t == `function` ? t(e) : t;
		}
		function Ro(e) {
			return zo(V(), R, e);
		}
		function zo(e, t, n) {
			var r = e.queue;
			if (r === null) throw Error(s(311));
			r.lastRenderedReducer = n;
			var i = e.baseQueue,
				a = r.pending;
			if (a !== null) {
				if (i !== null) {
					var o = i.next;
					(i.next = a.next), (a.next = o);
				}
				(t.baseQueue = i = a), (r.pending = null);
			}
			if (((a = e.baseState), i === null)) e.memoizedState = a;
			else {
				t = i.next;
				var c = (o = null),
					l = null,
					u = t,
					d = !1;
				do {
					var f = u.lane & -536870913;
					if (f === u.lane ? (go & f) === f : (J & f) === f) {
						var p = u.revertLane;
						if (p === 0)
							l !== null &&
								(l = l.next =
									{
										lane: 0,
										revertLane: 0,
										gesture: null,
										action: u.action,
										hasEagerState: u.hasEagerState,
										eagerState: u.eagerState,
										next: null,
									}),
								f === ga && (d = !0);
						else if ((go & p) === p) {
							(u = u.next), p === ga && (d = !0);
							continue;
						} else
							(f = {
								lane: 0,
								revertLane: u.revertLane,
								gesture: null,
								action: u.action,
								hasEagerState: u.hasEagerState,
								eagerState: u.eagerState,
								next: null,
							}),
								l === null ? ((c = l = f), (o = a)) : (l = l.next = f),
								(L.lanes |= p),
								(Gl |= p);
						(f = u.action), yo && n(a, f), (a = u.hasEagerState ? u.eagerState : n(a, f));
					} else
						(p = {
							lane: f,
							revertLane: u.revertLane,
							gesture: u.gesture,
							action: u.action,
							hasEagerState: u.hasEagerState,
							eagerState: u.eagerState,
							next: null,
						}),
							l === null ? ((c = l = p), (o = a)) : (l = l.next = p),
							(L.lanes |= f),
							(Gl |= f);
					u = u.next;
				} while (u !== null && u !== t);
				if (
					(l === null ? (o = a) : (l.next = c),
					!Or(a, e.memoizedState) && ((H = !0), d && ((n = _a), n !== null)))
				)
					throw n;
				(e.memoizedState = a), (e.baseState = o), (e.baseQueue = l), (r.lastRenderedState = a);
			}
			return i === null && (r.lanes = 0), [e.memoizedState, r.dispatch];
		}
		function Bo(e) {
			var t = V(),
				n = t.queue;
			if (n === null) throw Error(s(311));
			n.lastRenderedReducer = e;
			var r = n.dispatch,
				i = n.pending,
				a = t.memoizedState;
			if (i !== null) {
				n.pending = null;
				var o = (i = i.next);
				do (a = e(a, o.action)), (o = o.next);
				while (o !== i);
				Or(a, t.memoizedState) || (H = !0),
					(t.memoizedState = a),
					t.baseQueue === null && (t.baseState = a),
					(n.lastRenderedState = a);
			}
			return [a, r];
		}
		function Vo(e, t, n) {
			var r = L,
				i = V(),
				a = P;
			if (a) {
				if (n === void 0) throw Error(s(407));
				n = n();
			} else n = t();
			var o = !Or((R || i).memoizedState, n);
			if (
				(o && ((i.memoizedState = n), (H = !0)),
				(i = i.queue),
				fs(Wo.bind(null, r, i, e), [e]),
				i.getSnapshot !== t || o || (z !== null && z.memoizedState.tag & 1))
			) {
				if (
					((r.flags |= 2048),
					ss(9, { destroy: void 0 }, Uo.bind(null, r, i, n, t), null),
					K === null)
				)
					throw Error(s(349));
				a || go & 127 || Ho(r, t, n);
			}
			return n;
		}
		function Ho(e, t, n) {
			(e.flags |= 16384),
				(e = { getSnapshot: t, value: n }),
				(t = L.updateQueue),
				t === null
					? ((t = No()), (L.updateQueue = t), (t.stores = [e]))
					: ((n = t.stores), n === null ? (t.stores = [e]) : n.push(e));
		}
		function Uo(e, t, n, r) {
			(t.value = n), (t.getSnapshot = r), Go(t) && Ko(e);
		}
		function Wo(e, t, n) {
			return n(function () {
				Go(t) && Ko(e);
			});
		}
		function Go(e) {
			var t = e.getSnapshot;
			e = e.value;
			try {
				var n = t();
				return !Or(e, n);
			} catch {
				return !0;
			}
		}
		function Ko(e) {
			var t = li(e, 2);
			t !== null && hu(t, e, 2);
		}
		function qo(e) {
			var t = Mo();
			if (typeof e == `function`) {
				var n = e;
				if (((e = n()), yo)) {
					We(!0);
					try {
						n();
					} finally {
						We(!1);
					}
				}
			}
			return (
				(t.memoizedState = t.baseState = e),
				(t.queue = {
					pending: null,
					lanes: 0,
					dispatch: null,
					lastRenderedReducer: Lo,
					lastRenderedState: e,
				}),
				t
			);
		}
		function Jo(e, t, n, r) {
			return (e.baseState = n), zo(e, R, typeof r == `function` ? r : Lo);
		}
		function Yo(e, t, n, r, i) {
			if (Ls(e)) throw Error(s(485));
			if (((e = t.action), e !== null)) {
				var a = {
					payload: i,
					action: e,
					next: null,
					isTransition: !0,
					status: `pending`,
					value: null,
					reason: null,
					listeners: [],
					then: function (e) {
						a.listeners.push(e);
					},
				};
				T.T === null ? (a.isTransition = !1) : n(!0),
					r(a),
					(n = t.pending),
					n === null
						? ((a.next = t.pending = a), Xo(t, a))
						: ((a.next = n.next), (t.pending = n.next = a));
			}
		}
		function Xo(e, t) {
			var n = t.action,
				r = t.payload,
				i = e.state;
			if (t.isTransition) {
				var a = T.T,
					o = {};
				T.T = o;
				try {
					var s = n(i, r),
						c = T.S;
					c !== null && c(o, s), Zo(e, t, s);
				} catch (n) {
					$o(e, t, n);
				} finally {
					a !== null && o.types !== null && (a.types = o.types), (T.T = a);
				}
			} else
				try {
					(a = n(i, r)), Zo(e, t, a);
				} catch (n) {
					$o(e, t, n);
				}
		}
		function Zo(e, t, n) {
			typeof n == `object` && n && typeof n.then == `function`
				? n.then(
						function (n) {
							Qo(e, t, n);
						},
						function (n) {
							return $o(e, t, n);
						},
					)
				: Qo(e, t, n);
		}
		function Qo(e, t, n) {
			(t.status = `fulfilled`),
				(t.value = n),
				es(t),
				(e.state = n),
				(t = e.pending),
				t !== null &&
					((n = t.next), n === t ? (e.pending = null) : ((n = n.next), (t.next = n), Xo(e, n)));
		}
		function $o(e, t, n) {
			var r = e.pending;
			if (((e.pending = null), r !== null)) {
				r = r.next;
				do (t.status = `rejected`), (t.reason = n), es(t), (t = t.next);
				while (t !== r);
			}
			e.action = null;
		}
		function es(e) {
			e = e.listeners;
			for (var t = 0; t < e.length; t++) (0, e[t])();
		}
		function ts(e, t) {
			return t;
		}
		function ns(e, t) {
			if (P) {
				var n = K.formState;
				if (n !== null) {
					a: {
						var r = L;
						if (P) {
							if (N) {
								b: {
									for (var i = N, a = Vi; i.nodeType !== 8; ) {
										if (!a) {
											i = null;
											break b;
										}
										if (((i = cf(i.nextSibling)), i === null)) {
											i = null;
											break b;
										}
									}
									(a = i.data), (i = a === `F!` || a === `F` ? i : null);
								}
								if (i) {
									(N = cf(i.nextSibling)), (r = i.data === `F!`);
									break a;
								}
							}
							Ui(r);
						}
						r = !1;
					}
					r && (t = n[0]);
				}
			}
			return (
				(n = Mo()),
				(n.memoizedState = n.baseState = t),
				(r = {
					pending: null,
					lanes: 0,
					dispatch: null,
					lastRenderedReducer: ts,
					lastRenderedState: t,
				}),
				(n.queue = r),
				(n = Ps.bind(null, L, r)),
				(r.dispatch = n),
				(r = qo(!1)),
				(a = Is.bind(null, L, !1, r.queue)),
				(r = Mo()),
				(i = { state: t, dispatch: null, action: e, pending: null }),
				(r.queue = i),
				(n = Yo.bind(null, L, i, a, n)),
				(i.dispatch = n),
				(r.memoizedState = e),
				[t, n, !1]
			);
		}
		function rs(e) {
			return is(V(), R, e);
		}
		function is(e, t, n) {
			if (
				((t = zo(e, t, ts)[0]),
				(e = Ro(Lo)[0]),
				typeof t == `object` && t && typeof t.then == `function`)
			)
				try {
					var r = Po(t);
				} catch (e) {
					throw e === Ea ? Oa : e;
				}
			else r = t;
			t = V();
			var i = t.queue,
				a = i.dispatch;
			return (
				n !== t.memoizedState &&
					((L.flags |= 2048), ss(9, { destroy: void 0 }, as.bind(null, i, n), null)),
				[r, a, e]
			);
		}
		function as(e, t) {
			e.action = t;
		}
		function os(e) {
			var t = V(),
				n = R;
			if (n !== null) return is(t, n, e);
			V(), (t = t.memoizedState), (n = V());
			var r = n.queue.dispatch;
			return (n.memoizedState = e), [t, r, !1];
		}
		function ss(e, t, n, r) {
			return (
				(e = { tag: e, create: n, deps: r, inst: t, next: null }),
				(t = L.updateQueue),
				t === null && ((t = No()), (L.updateQueue = t)),
				(n = t.lastEffect),
				n === null
					? (t.lastEffect = e.next = e)
					: ((r = n.next), (n.next = e), (e.next = r), (t.lastEffect = e)),
				e
			);
		}
		function cs() {
			return V().memoizedState;
		}
		function ls(e, t, n, r) {
			var i = Mo();
			(L.flags |= e),
				(i.memoizedState = ss(1 | t, { destroy: void 0 }, n, r === void 0 ? null : r));
		}
		function us(e, t, n, r) {
			var i = V();
			r = r === void 0 ? null : r;
			var a = i.memoizedState.inst;
			R !== null && r !== null && wo(r, R.memoizedState.deps)
				? (i.memoizedState = ss(t, a, n, r))
				: ((L.flags |= e), (i.memoizedState = ss(1 | t, a, n, r)));
		}
		function ds(e, t) {
			ls(8390656, 8, e, t);
		}
		function fs(e, t) {
			us(2048, 8, e, t);
		}
		function ps(e) {
			L.flags |= 4;
			var t = L.updateQueue;
			if (t === null) (t = No()), (L.updateQueue = t), (t.events = [e]);
			else {
				var n = t.events;
				n === null ? (t.events = [e]) : n.push(e);
			}
		}
		function ms(e) {
			var t = V().memoizedState;
			return (
				ps({ ref: t, nextImpl: e }),
				function () {
					if (G & 2) throw Error(s(440));
					return t.impl.apply(void 0, arguments);
				}
			);
		}
		function hs(e, t) {
			return us(4, 2, e, t);
		}
		function gs(e, t) {
			return us(4, 4, e, t);
		}
		function _s(e, t) {
			if (typeof t == `function`) {
				e = e();
				var n = t(e);
				return function () {
					typeof n == `function` ? n() : t(null);
				};
			}
			if (t != null)
				return (
					(e = e()),
					(t.current = e),
					function () {
						t.current = null;
					}
				);
		}
		function vs(e, t, n) {
			(n = n == null ? null : n.concat([e])), us(4, 4, _s.bind(null, t, e), n);
		}
		function ys() {}
		function bs(e, t) {
			var n = V();
			t = t === void 0 ? null : t;
			var r = n.memoizedState;
			return t !== null && wo(t, r[1]) ? r[0] : ((n.memoizedState = [e, t]), e);
		}
		function xs(e, t) {
			var n = V();
			t = t === void 0 ? null : t;
			var r = n.memoizedState;
			if (t !== null && wo(t, r[1])) return r[0];
			if (((r = e()), yo)) {
				We(!0);
				try {
					e();
				} finally {
					We(!1);
				}
			}
			return (n.memoizedState = [r, t]), r;
		}
		function Ss(e, t, n) {
			return n === void 0 || (go & 1073741824 && !(J & 261930))
				? (e.memoizedState = t)
				: ((e.memoizedState = n), (e = mu()), (L.lanes |= e), (Gl |= e), n);
		}
		function Cs(e, t, n, r) {
			return Or(n, t)
				? n
				: no.current === null
					? !(go & 42) || (go & 1073741824 && !(J & 261930))
						? ((H = !0), (e.memoizedState = n))
						: ((e = mu()), (L.lanes |= e), (Gl |= e), t)
					: ((e = Ss(e, n, r)), Or(e, t) || (H = !0), e);
		}
		function ws(e, t, n, r, i) {
			var a = E.p;
			E.p = a !== 0 && 8 > a ? a : 8;
			var o = T.T,
				s = {};
			(T.T = s), Is(e, !1, t, n);
			try {
				var c = i(),
					l = T.S;
				l !== null && l(s, c),
					typeof c == `object` && c && typeof c.then == `function`
						? Fs(e, t, ba(c, r), pu(e))
						: Fs(e, t, r, pu(e));
			} catch (n) {
				Fs(e, t, { then: function () {}, status: `rejected`, reason: n }, pu());
			} finally {
				(E.p = a), o !== null && s.types !== null && (o.types = s.types), (T.T = o);
			}
		}
		function Ts() {}
		function Es(e, t, n, r) {
			if (e.tag !== 5) throw Error(s(476));
			var i = Ds(e).queue;
			ws(
				e,
				i,
				t,
				de,
				n === null
					? Ts
					: function () {
							return Os(e), n(r);
						},
			);
		}
		function Ds(e) {
			var t = e.memoizedState;
			if (t !== null) return t;
			t = {
				memoizedState: de,
				baseState: de,
				baseQueue: null,
				queue: {
					pending: null,
					lanes: 0,
					dispatch: null,
					lastRenderedReducer: Lo,
					lastRenderedState: de,
				},
				next: null,
			};
			var n = {};
			return (
				(t.next = {
					memoizedState: n,
					baseState: n,
					baseQueue: null,
					queue: {
						pending: null,
						lanes: 0,
						dispatch: null,
						lastRenderedReducer: Lo,
						lastRenderedState: n,
					},
					next: null,
				}),
				(e.memoizedState = t),
				(e = e.alternate),
				e !== null && (e.memoizedState = t),
				t
			);
		}
		function Os(e) {
			var t = Ds(e);
			t.next === null && (t = e.alternate.memoizedState), Fs(e, t.next.queue, {}, pu());
		}
		function ks() {
			return oa(Qf);
		}
		function As() {
			return V().memoizedState;
		}
		function js() {
			return V().memoizedState;
		}
		function Ms(e) {
			for (var t = e.return; t !== null; ) {
				switch (t.tag) {
					case 24:
					case 3:
						var n = pu();
						e = qa(n);
						var r = Ja(t, e, n);
						r !== null && (hu(r, t, n), Ya(r, t, n)), (t = { cache: fa() }), (e.payload = t);
						return;
				}
				t = t.return;
			}
		}
		function Ns(e, t, n) {
			var r = pu();
			(n = {
				lane: r,
				revertLane: 0,
				gesture: null,
				action: n,
				hasEagerState: !1,
				eagerState: null,
				next: null,
			}),
				Ls(e) ? Rs(t, n) : ((n = ci(e, t, n, r)), n !== null && (hu(n, e, r), zs(n, t, r)));
		}
		function Ps(e, t, n) {
			Fs(e, t, n, pu());
		}
		function Fs(e, t, n, r) {
			var i = {
				lane: r,
				revertLane: 0,
				gesture: null,
				action: n,
				hasEagerState: !1,
				eagerState: null,
				next: null,
			};
			if (Ls(e)) Rs(t, i);
			else {
				var a = e.alternate;
				if (
					e.lanes === 0 &&
					(a === null || a.lanes === 0) &&
					((a = t.lastRenderedReducer), a !== null)
				)
					try {
						var o = t.lastRenderedState,
							s = a(o, n);
						if (((i.hasEagerState = !0), (i.eagerState = s), Or(s, o)))
							return si(e, t, i, 0), K === null && oi(), !1;
					} catch {}
				if (((n = ci(e, t, i, r)), n !== null)) return hu(n, e, r), zs(n, t, r), !0;
			}
			return !1;
		}
		function Is(e, t, n, r) {
			if (
				((r = {
					lane: 2,
					revertLane: dd(),
					gesture: null,
					action: r,
					hasEagerState: !1,
					eagerState: null,
					next: null,
				}),
				Ls(e))
			) {
				if (t) throw Error(s(479));
			} else (t = ci(e, n, r, 2)), t !== null && hu(t, e, 2);
		}
		function Ls(e) {
			var t = e.alternate;
			return e === L || (t !== null && t === L);
		}
		function Rs(e, t) {
			vo = _o = !0;
			var n = e.pending;
			n === null ? (t.next = t) : ((t.next = n.next), (n.next = t)), (e.pending = t);
		}
		function zs(e, t, n) {
			if (n & 4194048) {
				var r = t.lanes;
				(r &= e.pendingLanes), (n |= r), (t.lanes = n), st(e, n);
			}
		}
		var Bs = {
			readContext: oa,
			use: Fo,
			useCallback: B,
			useContext: B,
			useEffect: B,
			useImperativeHandle: B,
			useLayoutEffect: B,
			useInsertionEffect: B,
			useMemo: B,
			useReducer: B,
			useRef: B,
			useState: B,
			useDebugValue: B,
			useDeferredValue: B,
			useTransition: B,
			useSyncExternalStore: B,
			useId: B,
			useHostTransitionStatus: B,
			useFormState: B,
			useActionState: B,
			useOptimistic: B,
			useMemoCache: B,
			useCacheRefresh: B,
		};
		Bs.useEffectEvent = B;
		var Vs = {
				readContext: oa,
				use: Fo,
				useCallback: function (e, t) {
					return (Mo().memoizedState = [e, t === void 0 ? null : t]), e;
				},
				useContext: oa,
				useEffect: ds,
				useImperativeHandle: function (e, t, n) {
					(n = n == null ? null : n.concat([e])), ls(4194308, 4, _s.bind(null, t, e), n);
				},
				useLayoutEffect: function (e, t) {
					return ls(4194308, 4, e, t);
				},
				useInsertionEffect: function (e, t) {
					ls(4, 2, e, t);
				},
				useMemo: function (e, t) {
					var n = Mo();
					t = t === void 0 ? null : t;
					var r = e();
					if (yo) {
						We(!0);
						try {
							e();
						} finally {
							We(!1);
						}
					}
					return (n.memoizedState = [r, t]), r;
				},
				useReducer: function (e, t, n) {
					var r = Mo();
					if (n !== void 0) {
						var i = n(t);
						if (yo) {
							We(!0);
							try {
								n(t);
							} finally {
								We(!1);
							}
						}
					} else i = t;
					return (
						(r.memoizedState = r.baseState = i),
						(e = {
							pending: null,
							lanes: 0,
							dispatch: null,
							lastRenderedReducer: e,
							lastRenderedState: i,
						}),
						(r.queue = e),
						(e = e.dispatch = Ns.bind(null, L, e)),
						[r.memoizedState, e]
					);
				},
				useRef: function (e) {
					var t = Mo();
					return (e = { current: e }), (t.memoizedState = e);
				},
				useState: function (e) {
					e = qo(e);
					var t = e.queue,
						n = Ps.bind(null, L, t);
					return (t.dispatch = n), [e.memoizedState, n];
				},
				useDebugValue: ys,
				useDeferredValue: function (e, t) {
					return Ss(Mo(), e, t);
				},
				useTransition: function () {
					var e = qo(!1);
					return (e = ws.bind(null, L, e.queue, !0, !1)), (Mo().memoizedState = e), [!1, e];
				},
				useSyncExternalStore: function (e, t, n) {
					var r = L,
						i = Mo();
					if (P) {
						if (n === void 0) throw Error(s(407));
						n = n();
					} else {
						if (((n = t()), K === null)) throw Error(s(349));
						J & 127 || Ho(r, t, n);
					}
					i.memoizedState = n;
					var a = { value: n, getSnapshot: t };
					return (
						(i.queue = a),
						ds(Wo.bind(null, r, a, e), [e]),
						(r.flags |= 2048),
						ss(9, { destroy: void 0 }, Uo.bind(null, r, a, n, t), null),
						n
					);
				},
				useId: function () {
					var e = Mo(),
						t = K.identifierPrefix;
					if (P) {
						var n = Ni,
							r = Mi;
						(n = (r & ~(1 << (32 - Ge(r) - 1))).toString(32) + n),
							(t = `_` + t + `R_` + n),
							(n = bo++),
							0 < n && (t += `H` + n.toString(32)),
							(t += `_`);
					} else (n = Co++), (t = `_` + t + `r_` + n.toString(32) + `_`);
					return (e.memoizedState = t);
				},
				useHostTransitionStatus: ks,
				useFormState: ns,
				useActionState: ns,
				useOptimistic: function (e) {
					var t = Mo();
					t.memoizedState = t.baseState = e;
					var n = {
						pending: null,
						lanes: 0,
						dispatch: null,
						lastRenderedReducer: null,
						lastRenderedState: null,
					};
					return (t.queue = n), (t = Is.bind(null, L, !0, n)), (n.dispatch = t), [e, t];
				},
				useMemoCache: Io,
				useCacheRefresh: function () {
					return (Mo().memoizedState = Ms.bind(null, L));
				},
				useEffectEvent: function (e) {
					var t = Mo(),
						n = { impl: e };
					return (
						(t.memoizedState = n),
						function () {
							if (G & 2) throw Error(s(440));
							return n.impl.apply(void 0, arguments);
						}
					);
				},
			},
			Hs = {
				readContext: oa,
				use: Fo,
				useCallback: bs,
				useContext: oa,
				useEffect: fs,
				useImperativeHandle: vs,
				useInsertionEffect: hs,
				useLayoutEffect: gs,
				useMemo: xs,
				useReducer: Ro,
				useRef: cs,
				useState: function () {
					return Ro(Lo);
				},
				useDebugValue: ys,
				useDeferredValue: function (e, t) {
					return Cs(V(), R.memoizedState, e, t);
				},
				useTransition: function () {
					var e = Ro(Lo)[0],
						t = V().memoizedState;
					return [typeof e == `boolean` ? e : Po(e), t];
				},
				useSyncExternalStore: Vo,
				useId: As,
				useHostTransitionStatus: ks,
				useFormState: rs,
				useActionState: rs,
				useOptimistic: function (e, t) {
					return Jo(V(), R, e, t);
				},
				useMemoCache: Io,
				useCacheRefresh: js,
			};
		Hs.useEffectEvent = ms;
		var Us = {
			readContext: oa,
			use: Fo,
			useCallback: bs,
			useContext: oa,
			useEffect: fs,
			useImperativeHandle: vs,
			useInsertionEffect: hs,
			useLayoutEffect: gs,
			useMemo: xs,
			useReducer: Bo,
			useRef: cs,
			useState: function () {
				return Bo(Lo);
			},
			useDebugValue: ys,
			useDeferredValue: function (e, t) {
				var n = V();
				return R === null ? Ss(n, e, t) : Cs(n, R.memoizedState, e, t);
			},
			useTransition: function () {
				var e = Bo(Lo)[0],
					t = V().memoizedState;
				return [typeof e == `boolean` ? e : Po(e), t];
			},
			useSyncExternalStore: Vo,
			useId: As,
			useHostTransitionStatus: ks,
			useFormState: os,
			useActionState: os,
			useOptimistic: function (e, t) {
				var n = V();
				return R === null ? ((n.baseState = e), [e, n.queue.dispatch]) : Jo(n, R, e, t);
			},
			useMemoCache: Io,
			useCacheRefresh: js,
		};
		Us.useEffectEvent = ms;
		function Ws(e, t, n, r) {
			(t = e.memoizedState),
				(n = n(r, t)),
				(n = n == null ? t : h({}, t, n)),
				(e.memoizedState = n),
				e.lanes === 0 && (e.updateQueue.baseState = n);
		}
		var Gs = {
			enqueueSetState: function (e, t, n) {
				e = e._reactInternals;
				var r = pu(),
					i = qa(r);
				(i.payload = t),
					n != null && (i.callback = n),
					(t = Ja(e, i, r)),
					t !== null && (hu(t, e, r), Ya(t, e, r));
			},
			enqueueReplaceState: function (e, t, n) {
				e = e._reactInternals;
				var r = pu(),
					i = qa(r);
				(i.tag = 1),
					(i.payload = t),
					n != null && (i.callback = n),
					(t = Ja(e, i, r)),
					t !== null && (hu(t, e, r), Ya(t, e, r));
			},
			enqueueForceUpdate: function (e, t) {
				e = e._reactInternals;
				var n = pu(),
					r = qa(n);
				(r.tag = 2),
					t != null && (r.callback = t),
					(t = Ja(e, r, n)),
					t !== null && (hu(t, e, n), Ya(t, e, n));
			},
		};
		function Ks(e, t, n, r, i, a, o) {
			return (
				(e = e.stateNode),
				typeof e.shouldComponentUpdate == `function`
					? e.shouldComponentUpdate(r, a, o)
					: t.prototype && t.prototype.isPureReactComponent
						? !kr(n, r) || !kr(i, a)
						: !0
			);
		}
		function qs(e, t, n, r) {
			(e = t.state),
				typeof t.componentWillReceiveProps == `function` && t.componentWillReceiveProps(n, r),
				typeof t.UNSAFE_componentWillReceiveProps == `function` &&
					t.UNSAFE_componentWillReceiveProps(n, r),
				t.state !== e && Gs.enqueueReplaceState(t, t.state, null);
		}
		function Js(e, t) {
			var n = t;
			if (`ref` in t) for (var r in ((n = {}), t)) r !== `ref` && (n[r] = t[r]);
			if ((e = e.defaultProps))
				for (var i in (n === t && (n = h({}, n)), e)) n[i] === void 0 && (n[i] = e[i]);
			return n;
		}
		function Ys(e) {
			ni(e);
		}
		function Xs(e) {
			console.error(e);
		}
		function Zs(e) {
			ni(e);
		}
		function Qs(e, t) {
			try {
				var n = e.onUncaughtError;
				n(t.value, { componentStack: t.stack });
			} catch (e) {
				setTimeout(function () {
					throw e;
				});
			}
		}
		function $s(e, t, n) {
			try {
				var r = e.onCaughtError;
				r(n.value, { componentStack: n.stack, errorBoundary: t.tag === 1 ? t.stateNode : null });
			} catch (e) {
				setTimeout(function () {
					throw e;
				});
			}
		}
		function ec(e, t, n) {
			return (
				(n = qa(n)),
				(n.tag = 3),
				(n.payload = { element: null }),
				(n.callback = function () {
					Qs(e, t);
				}),
				n
			);
		}
		function tc(e) {
			return (e = qa(e)), (e.tag = 3), e;
		}
		function nc(e, t, n, r) {
			var i = n.type.getDerivedStateFromError;
			if (typeof i == `function`) {
				var a = r.value;
				(e.payload = function () {
					return i(a);
				}),
					(e.callback = function () {
						$s(t, n, r);
					});
			}
			var o = n.stateNode;
			o !== null &&
				typeof o.componentDidCatch == `function` &&
				(e.callback = function () {
					$s(t, n, r),
						typeof i != `function` && (ru === null ? (ru = new Set([this])) : ru.add(this));
					var e = r.stack;
					this.componentDidCatch(r.value, { componentStack: e === null ? `` : e });
				});
		}
		function rc(e, t, n, r, i) {
			if (((n.flags |= 32768), typeof r == `object` && r && typeof r.then == `function`)) {
				if (((t = n.alternate), t !== null && ra(t, n, i, !0), (n = so.current), n !== null)) {
					switch (n.tag) {
						case 31:
						case 13:
							return (
								co === null ? Du() : n.alternate === null && X === 0 && (X = 3),
								(n.flags &= -257),
								(n.flags |= 65536),
								(n.lanes = i),
								r === ka
									? (n.flags |= 16384)
									: ((t = n.updateQueue),
										t === null ? (n.updateQueue = new Set([r])) : t.add(r),
										Gu(e, r, i)),
								!1
							);
						case 22:
							return (
								(n.flags |= 65536),
								r === ka
									? (n.flags |= 16384)
									: ((t = n.updateQueue),
										t === null
											? ((t = {
													transitions: null,
													markerInstances: null,
													retryQueue: new Set([r]),
												}),
												(n.updateQueue = t))
											: ((n = t.retryQueue), n === null ? (t.retryQueue = new Set([r])) : n.add(r)),
										Gu(e, r, i)),
								!1
							);
					}
					throw Error(s(435, n.tag));
				}
				return Gu(e, r, i), Du(), !1;
			}
			if (P)
				return (
					(t = so.current),
					t === null
						? (r !== Hi && ((t = Error(s(423), { cause: r })), Yi(wi(t, n))),
							(e = e.current.alternate),
							(e.flags |= 65536),
							(i &= -i),
							(e.lanes |= i),
							(r = wi(r, n)),
							(i = ec(e.stateNode, r, i)),
							Xa(e, i),
							X !== 4 && (X = 2))
						: (!(t.flags & 65536) && (t.flags |= 256),
							(t.flags |= 65536),
							(t.lanes = i),
							r !== Hi && ((e = Error(s(422), { cause: r })), Yi(wi(e, n)))),
					!1
				);
			var a = Error(s(520), { cause: r });
			if (((a = wi(a, n)), Xl === null ? (Xl = [a]) : Xl.push(a), X !== 4 && (X = 2), t === null))
				return !0;
			(r = wi(r, n)), (n = t);
			do {
				switch (n.tag) {
					case 3:
						return (
							(n.flags |= 65536),
							(e = i & -i),
							(n.lanes |= e),
							(e = ec(n.stateNode, r, e)),
							Xa(n, e),
							!1
						);
					case 1:
						if (
							((t = n.type),
							(a = n.stateNode),
							!(n.flags & 128) &&
								(typeof t.getDerivedStateFromError == `function` ||
									(a !== null &&
										typeof a.componentDidCatch == `function` &&
										(ru === null || !ru.has(a)))))
						)
							return (
								(n.flags |= 65536),
								(i &= -i),
								(n.lanes |= i),
								(i = tc(i)),
								nc(i, e, n, r),
								Xa(n, i),
								!1
							);
				}
				n = n.return;
			} while (n !== null);
			return !1;
		}
		var ic = Error(s(461)),
			H = !1;
		function ac(e, t, n, r) {
			t.child = e === null ? Ua(t, null, n, r) : Ha(t, e.child, n, r);
		}
		function oc(e, t, n, r, i) {
			n = n.render;
			var a = t.ref;
			if (`ref` in r) {
				var o = {};
				for (var s in r) s !== `ref` && (o[s] = r[s]);
			} else o = r;
			return (
				aa(t),
				(r = To(e, t, n, o, a, i)),
				(s = ko()),
				e !== null && !H
					? (Ao(e, t, i), Ac(e, t, i))
					: (P && s && Ii(t), (t.flags |= 1), ac(e, t, r, i), t.child)
			);
		}
		function sc(e, t, n, r, i) {
			if (e === null) {
				var a = n.type;
				return typeof a == `function` && !hi(a) && a.defaultProps === void 0 && n.compare === null
					? ((t.tag = 15), (t.type = a), cc(e, t, a, r, i))
					: ((e = vi(n.type, null, r, t, t.mode, i)),
						(e.ref = t.ref),
						(e.return = t),
						(t.child = e));
			}
			if (((a = e.child), !jc(e, i))) {
				var o = a.memoizedProps;
				if (((n = n.compare), (n = n === null ? kr : n), n(o, r) && e.ref === t.ref))
					return Ac(e, t, i);
			}
			return (t.flags |= 1), (e = gi(a, r)), (e.ref = t.ref), (e.return = t), (t.child = e);
		}
		function cc(e, t, n, r, i) {
			if (e !== null) {
				var a = e.memoizedProps;
				if (kr(a, r) && e.ref === t.ref)
					if (((H = !1), (t.pendingProps = r = a), jc(e, i))) e.flags & 131072 && (H = !0);
					else return (t.lanes = e.lanes), Ac(e, t, i);
			}
			return gc(e, t, n, r, i);
		}
		function lc(e, t, n, r) {
			var i = r.children,
				a = e === null ? null : e.memoizedState;
			if (
				(e === null &&
					t.stateNode === null &&
					(t.stateNode = {
						_visibility: 1,
						_pendingMarkers: null,
						_retryCache: null,
						_transitions: null,
					}),
				r.mode === `hidden`)
			) {
				if (t.flags & 128) {
					if (((a = a === null ? n : a.baseLanes | n), e !== null)) {
						for (r = t.child = e.child, i = 0; r !== null; )
							(i = i | r.lanes | r.childLanes), (r = r.sibling);
						r = i & ~a;
					} else (r = 0), (t.child = null);
					return dc(e, t, a, n, r);
				}
				if (n & 536870912)
					(t.memoizedState = { baseLanes: 0, cachePool: null }),
						e !== null && wa(t, a === null ? null : a.cachePool),
						a === null ? ao() : io(t, a),
						fo(t);
				else return (r = t.lanes = 536870912), dc(e, t, a === null ? n : a.baseLanes | n, n, r);
			} else
				a === null
					? (e !== null && wa(t, null), ao(), po(t))
					: (wa(t, a.cachePool), io(t, a), po(t), (t.memoizedState = null));
			return ac(e, t, i, n), t.child;
		}
		function uc(e, t) {
			return (
				(e !== null && e.tag === 22) ||
					t.stateNode !== null ||
					(t.stateNode = {
						_visibility: 1,
						_pendingMarkers: null,
						_retryCache: null,
						_transitions: null,
					}),
				t.sibling
			);
		}
		function dc(e, t, n, r, i) {
			var a = Ca();
			return (
				(a = a === null ? null : { parent: F._currentValue, pool: a }),
				(t.memoizedState = { baseLanes: n, cachePool: a }),
				e !== null && wa(t, null),
				ao(),
				fo(t),
				e !== null && ra(e, t, r, !0),
				(t.childLanes = i),
				null
			);
		}
		function fc(e, t) {
			return (
				(t = Tc({ mode: t.mode, children: t.children }, e.mode)),
				(t.ref = e.ref),
				(e.child = t),
				(t.return = e),
				t
			);
		}
		function pc(e, t, n) {
			return (
				Ha(t, e.child, null, n),
				(e = fc(t, t.pendingProps)),
				(e.flags |= 2),
				mo(t),
				(t.memoizedState = null),
				e
			);
		}
		function mc(e, t, n) {
			var r = t.pendingProps,
				i = (t.flags & 128) != 0;
			if (((t.flags &= -129), e === null)) {
				if (P) {
					if (r.mode === `hidden`) return (e = fc(t, r)), (t.lanes = 536870912), uc(null, e);
					if (
						(uo(t),
						(e = N)
							? ((e = rf(e, Vi)),
								(e = e !== null && e.data === `&` ? e : null),
								e !== null &&
									((t.memoizedState = {
										dehydrated: e,
										treeContext: ji === null ? null : { id: Mi, overflow: Ni },
										retryLane: 536870912,
										hydrationErrors: null,
									}),
									(n = xi(e)),
									(n.return = t),
									(t.child = n),
									(zi = t),
									(N = null)))
							: (e = null),
						e === null)
					)
						throw Ui(t);
					return (t.lanes = 536870912), null;
				}
				return fc(t, r);
			}
			var a = e.memoizedState;
			if (a !== null) {
				var o = a.dehydrated;
				if ((uo(t), i))
					if (t.flags & 256) (t.flags &= -257), (t = pc(e, t, n));
					else if (t.memoizedState !== null) (t.child = e.child), (t.flags |= 128), (t = null);
					else throw Error(s(558));
				else if ((H || ra(e, t, n, !1), (i = (n & e.childLanes) !== 0), H || i)) {
					if (((r = K), r !== null && ((o = ct(r, n)), o !== 0 && o !== a.retryLane)))
						throw ((a.retryLane = o), li(e, o), hu(r, e, o), ic);
					Du(), (t = pc(e, t, n));
				} else
					(e = a.treeContext),
						(N = cf(o.nextSibling)),
						(zi = t),
						(P = !0),
						(Bi = null),
						(Vi = !1),
						e !== null && Ri(t, e),
						(t = fc(t, r)),
						(t.flags |= 4096);
				return t;
			}
			return (
				(e = gi(e.child, { mode: r.mode, children: r.children })),
				(e.ref = t.ref),
				(t.child = e),
				(e.return = t),
				e
			);
		}
		function hc(e, t) {
			var n = t.ref;
			if (n === null) e !== null && e.ref !== null && (t.flags |= 4194816);
			else {
				if (typeof n != `function` && typeof n != `object`) throw Error(s(284));
				(e === null || e.ref !== n) && (t.flags |= 4194816);
			}
		}
		function gc(e, t, n, r, i) {
			return (
				aa(t),
				(n = To(e, t, n, r, void 0, i)),
				(r = ko()),
				e !== null && !H
					? (Ao(e, t, i), Ac(e, t, i))
					: (P && r && Ii(t), (t.flags |= 1), ac(e, t, n, i), t.child)
			);
		}
		function _c(e, t, n, r, i, a) {
			return (
				aa(t),
				(t.updateQueue = null),
				(n = Do(t, r, n, i)),
				Eo(e),
				(r = ko()),
				e !== null && !H
					? (Ao(e, t, a), Ac(e, t, a))
					: (P && r && Ii(t), (t.flags |= 1), ac(e, t, n, a), t.child)
			);
		}
		function vc(e, t, n, r, i) {
			if ((aa(t), t.stateNode === null)) {
				var a = fi,
					o = n.contextType;
				typeof o == `object` && o && (a = oa(o)),
					(a = new n(r, a)),
					(t.memoizedState = a.state !== null && a.state !== void 0 ? a.state : null),
					(a.updater = Gs),
					(t.stateNode = a),
					(a._reactInternals = t),
					(a = t.stateNode),
					(a.props = r),
					(a.state = t.memoizedState),
					(a.refs = {}),
					Ga(t),
					(o = n.contextType),
					(a.context = typeof o == `object` && o ? oa(o) : fi),
					(a.state = t.memoizedState),
					(o = n.getDerivedStateFromProps),
					typeof o == `function` && (Ws(t, n, o, r), (a.state = t.memoizedState)),
					typeof n.getDerivedStateFromProps == `function` ||
						typeof a.getSnapshotBeforeUpdate == `function` ||
						(typeof a.UNSAFE_componentWillMount != `function` &&
							typeof a.componentWillMount != `function`) ||
						((o = a.state),
						typeof a.componentWillMount == `function` && a.componentWillMount(),
						typeof a.UNSAFE_componentWillMount == `function` && a.UNSAFE_componentWillMount(),
						o !== a.state && Gs.enqueueReplaceState(a, a.state, null),
						$a(t, r, a, i),
						Qa(),
						(a.state = t.memoizedState)),
					typeof a.componentDidMount == `function` && (t.flags |= 4194308),
					(r = !0);
			} else if (e === null) {
				a = t.stateNode;
				var s = t.memoizedProps,
					c = Js(n, s);
				a.props = c;
				var l = a.context,
					u = n.contextType;
				(o = fi), typeof u == `object` && u && (o = oa(u));
				var d = n.getDerivedStateFromProps;
				(u = typeof d == `function` || typeof a.getSnapshotBeforeUpdate == `function`),
					(s = t.pendingProps !== s),
					u ||
						(typeof a.UNSAFE_componentWillReceiveProps != `function` &&
							typeof a.componentWillReceiveProps != `function`) ||
						((s || l !== o) && qs(t, a, r, o)),
					(Wa = !1);
				var f = t.memoizedState;
				(a.state = f),
					$a(t, r, a, i),
					Qa(),
					(l = t.memoizedState),
					s || f !== l || Wa
						? (typeof d == `function` && (Ws(t, n, d, r), (l = t.memoizedState)),
							(c = Wa || Ks(t, n, c, r, f, l, o))
								? (u ||
										(typeof a.UNSAFE_componentWillMount != `function` &&
											typeof a.componentWillMount != `function`) ||
										(typeof a.componentWillMount == `function` && a.componentWillMount(),
										typeof a.UNSAFE_componentWillMount == `function` &&
											a.UNSAFE_componentWillMount()),
									typeof a.componentDidMount == `function` && (t.flags |= 4194308))
								: (typeof a.componentDidMount == `function` && (t.flags |= 4194308),
									(t.memoizedProps = r),
									(t.memoizedState = l)),
							(a.props = r),
							(a.state = l),
							(a.context = o),
							(r = c))
						: (typeof a.componentDidMount == `function` && (t.flags |= 4194308), (r = !1));
			} else {
				(a = t.stateNode),
					Ka(e, t),
					(o = t.memoizedProps),
					(u = Js(n, o)),
					(a.props = u),
					(d = t.pendingProps),
					(f = a.context),
					(l = n.contextType),
					(c = fi),
					typeof l == `object` && l && (c = oa(l)),
					(s = n.getDerivedStateFromProps),
					(l = typeof s == `function` || typeof a.getSnapshotBeforeUpdate == `function`) ||
						(typeof a.UNSAFE_componentWillReceiveProps != `function` &&
							typeof a.componentWillReceiveProps != `function`) ||
						((o !== d || f !== c) && qs(t, a, r, c)),
					(Wa = !1),
					(f = t.memoizedState),
					(a.state = f),
					$a(t, r, a, i),
					Qa();
				var p = t.memoizedState;
				o !== d || f !== p || Wa || (e !== null && e.dependencies !== null && ia(e.dependencies))
					? (typeof s == `function` && (Ws(t, n, s, r), (p = t.memoizedState)),
						(u =
							Wa ||
							Ks(t, n, u, r, f, p, c) ||
							(e !== null && e.dependencies !== null && ia(e.dependencies)))
							? (l ||
									(typeof a.UNSAFE_componentWillUpdate != `function` &&
										typeof a.componentWillUpdate != `function`) ||
									(typeof a.componentWillUpdate == `function` && a.componentWillUpdate(r, p, c),
									typeof a.UNSAFE_componentWillUpdate == `function` &&
										a.UNSAFE_componentWillUpdate(r, p, c)),
								typeof a.componentDidUpdate == `function` && (t.flags |= 4),
								typeof a.getSnapshotBeforeUpdate == `function` && (t.flags |= 1024))
							: (typeof a.componentDidUpdate != `function` ||
									(o === e.memoizedProps && f === e.memoizedState) ||
									(t.flags |= 4),
								typeof a.getSnapshotBeforeUpdate != `function` ||
									(o === e.memoizedProps && f === e.memoizedState) ||
									(t.flags |= 1024),
								(t.memoizedProps = r),
								(t.memoizedState = p)),
						(a.props = r),
						(a.state = p),
						(a.context = c),
						(r = u))
					: (typeof a.componentDidUpdate != `function` ||
							(o === e.memoizedProps && f === e.memoizedState) ||
							(t.flags |= 4),
						typeof a.getSnapshotBeforeUpdate != `function` ||
							(o === e.memoizedProps && f === e.memoizedState) ||
							(t.flags |= 1024),
						(r = !1));
			}
			return (
				(a = r),
				hc(e, t),
				(r = (t.flags & 128) != 0),
				a || r
					? ((a = t.stateNode),
						(n = r && typeof n.getDerivedStateFromError != `function` ? null : a.render()),
						(t.flags |= 1),
						e !== null && r
							? ((t.child = Ha(t, e.child, null, i)), (t.child = Ha(t, null, n, i)))
							: ac(e, t, n, i),
						(t.memoizedState = a.state),
						(e = t.child))
					: (e = Ac(e, t, i)),
				e
			);
		}
		function yc(e, t, n, r) {
			return qi(), (t.flags |= 256), ac(e, t, n, r), t.child;
		}
		var bc = { dehydrated: null, treeContext: null, retryLane: 0, hydrationErrors: null };
		function xc(e) {
			return { baseLanes: e, cachePool: Ta() };
		}
		function Sc(e, t, n) {
			return (e = e === null ? 0 : e.childLanes & ~n), t && (e |= Jl), e;
		}
		function Cc(e, t, n) {
			var r = t.pendingProps,
				i = !1,
				a = (t.flags & 128) != 0,
				o;
			if (
				((o = a) || (o = e !== null && e.memoizedState === null ? !1 : (I.current & 2) != 0),
				o && ((i = !0), (t.flags &= -129)),
				(o = (t.flags & 32) != 0),
				(t.flags &= -33),
				e === null)
			) {
				if (P) {
					if (
						(i ? lo(t) : po(t),
						(e = N)
							? ((e = rf(e, Vi)),
								(e = e !== null && e.data !== `&` ? e : null),
								e !== null &&
									((t.memoizedState = {
										dehydrated: e,
										treeContext: ji === null ? null : { id: Mi, overflow: Ni },
										retryLane: 536870912,
										hydrationErrors: null,
									}),
									(n = xi(e)),
									(n.return = t),
									(t.child = n),
									(zi = t),
									(N = null)))
							: (e = null),
						e === null)
					)
						throw Ui(t);
					return of(e) ? (t.lanes = 32) : (t.lanes = 536870912), null;
				}
				var c = r.children;
				return (
					(r = r.fallback),
					i
						? (po(t),
							(i = t.mode),
							(c = Tc({ mode: `hidden`, children: c }, i)),
							(r = yi(r, i, n, null)),
							(c.return = t),
							(r.return = t),
							(c.sibling = r),
							(t.child = c),
							(r = t.child),
							(r.memoizedState = xc(n)),
							(r.childLanes = Sc(e, o, n)),
							(t.memoizedState = bc),
							uc(null, r))
						: (lo(t), wc(t, c))
				);
			}
			var l = e.memoizedState;
			if (l !== null && ((c = l.dehydrated), c !== null)) {
				if (a)
					t.flags & 256
						? (lo(t), (t.flags &= -257), (t = Ec(e, t, n)))
						: t.memoizedState === null
							? (po(t),
								(c = r.fallback),
								(i = t.mode),
								(r = Tc({ mode: `visible`, children: r.children }, i)),
								(c = yi(c, i, n, null)),
								(c.flags |= 2),
								(r.return = t),
								(c.return = t),
								(r.sibling = c),
								(t.child = r),
								Ha(t, e.child, null, n),
								(r = t.child),
								(r.memoizedState = xc(n)),
								(r.childLanes = Sc(e, o, n)),
								(t.memoizedState = bc),
								(t = uc(null, r)))
							: (po(t), (t.child = e.child), (t.flags |= 128), (t = null));
				else if ((lo(t), of(c))) {
					if (((o = c.nextSibling && c.nextSibling.dataset), o)) var u = o.dgst;
					(o = u),
						(r = Error(s(419))),
						(r.stack = ``),
						(r.digest = o),
						Yi({ value: r, source: null, stack: null }),
						(t = Ec(e, t, n));
				} else if ((H || ra(e, t, n, !1), (o = (n & e.childLanes) !== 0), H || o)) {
					if (((o = K), o !== null && ((r = ct(o, n)), r !== 0 && r !== l.retryLane)))
						throw ((l.retryLane = r), li(e, r), hu(o, e, r), ic);
					af(c) || Du(), (t = Ec(e, t, n));
				} else
					af(c)
						? ((t.flags |= 192), (t.child = e.child), (t = null))
						: ((e = l.treeContext),
							(N = cf(c.nextSibling)),
							(zi = t),
							(P = !0),
							(Bi = null),
							(Vi = !1),
							e !== null && Ri(t, e),
							(t = wc(t, r.children)),
							(t.flags |= 4096));
				return t;
			}
			return i
				? (po(t),
					(c = r.fallback),
					(i = t.mode),
					(l = e.child),
					(u = l.sibling),
					(r = gi(l, { mode: `hidden`, children: r.children })),
					(r.subtreeFlags = l.subtreeFlags & 65011712),
					u === null ? ((c = yi(c, i, n, null)), (c.flags |= 2)) : (c = gi(u, c)),
					(c.return = t),
					(r.return = t),
					(r.sibling = c),
					(t.child = r),
					uc(null, r),
					(r = t.child),
					(c = e.child.memoizedState),
					c === null
						? (c = xc(n))
						: ((i = c.cachePool),
							i === null
								? (i = Ta())
								: ((l = F._currentValue), (i = i.parent === l ? i : { parent: l, pool: l })),
							(c = { baseLanes: c.baseLanes | n, cachePool: i })),
					(r.memoizedState = c),
					(r.childLanes = Sc(e, o, n)),
					(t.memoizedState = bc),
					uc(e.child, r))
				: (lo(t),
					(n = e.child),
					(e = n.sibling),
					(n = gi(n, { mode: `visible`, children: r.children })),
					(n.return = t),
					(n.sibling = null),
					e !== null &&
						((o = t.deletions), o === null ? ((t.deletions = [e]), (t.flags |= 16)) : o.push(e)),
					(t.child = n),
					(t.memoizedState = null),
					n);
		}
		function wc(e, t) {
			return (t = Tc({ mode: `visible`, children: t }, e.mode)), (t.return = e), (e.child = t);
		}
		function Tc(e, t) {
			return (e = mi(22, e, null, t)), (e.lanes = 0), e;
		}
		function Ec(e, t, n) {
			return (
				Ha(t, e.child, null, n),
				(e = wc(t, t.pendingProps.children)),
				(e.flags |= 2),
				(t.memoizedState = null),
				e
			);
		}
		function Dc(e, t, n) {
			e.lanes |= t;
			var r = e.alternate;
			r !== null && (r.lanes |= t), ta(e.return, t, n);
		}
		function Oc(e, t, n, r, i, a) {
			var o = e.memoizedState;
			o === null
				? (e.memoizedState = {
						isBackwards: t,
						rendering: null,
						renderingStartTime: 0,
						last: r,
						tail: n,
						tailMode: i,
						treeForkCount: a,
					})
				: ((o.isBackwards = t),
					(o.rendering = null),
					(o.renderingStartTime = 0),
					(o.last = r),
					(o.tail = n),
					(o.tailMode = i),
					(o.treeForkCount = a));
		}
		function kc(e, t, n) {
			var r = t.pendingProps,
				i = r.revealOrder,
				a = r.tail;
			r = r.children;
			var o = I.current,
				s = (o & 2) != 0;
			if (
				(s ? ((o = (o & 1) | 2), (t.flags |= 128)) : (o &= 1),
				k(I, o),
				ac(e, t, r, n),
				(r = P ? Oi : 0),
				!s && e !== null && e.flags & 128)
			)
				a: for (e = t.child; e !== null; ) {
					if (e.tag === 13) e.memoizedState !== null && Dc(e, n, t);
					else if (e.tag === 19) Dc(e, n, t);
					else if (e.child !== null) {
						(e.child.return = e), (e = e.child);
						continue;
					}
					if (e === t) break a;
					for (; e.sibling === null; ) {
						if (e.return === null || e.return === t) break a;
						e = e.return;
					}
					(e.sibling.return = e.return), (e = e.sibling);
				}
			switch (i) {
				case `forwards`:
					for (n = t.child, i = null; n !== null; )
						(e = n.alternate), e !== null && ho(e) === null && (i = n), (n = n.sibling);
					(n = i),
						n === null ? ((i = t.child), (t.child = null)) : ((i = n.sibling), (n.sibling = null)),
						Oc(t, !1, i, n, a, r);
					break;
				case `backwards`:
				case `unstable_legacy-backwards`:
					for (n = null, i = t.child, t.child = null; i !== null; ) {
						if (((e = i.alternate), e !== null && ho(e) === null)) {
							t.child = i;
							break;
						}
						(e = i.sibling), (i.sibling = n), (n = i), (i = e);
					}
					Oc(t, !0, n, null, a, r);
					break;
				case `together`:
					Oc(t, !1, null, null, void 0, r);
					break;
				default:
					t.memoizedState = null;
			}
			return t.child;
		}
		function Ac(e, t, n) {
			if (
				(e !== null && (t.dependencies = e.dependencies), (Gl |= t.lanes), (n & t.childLanes) === 0)
			)
				if (e !== null) {
					if ((ra(e, t, n, !1), (n & t.childLanes) === 0)) return null;
				} else return null;
			if (e !== null && t.child !== e.child) throw Error(s(153));
			if (t.child !== null) {
				for (
					e = t.child, n = gi(e, e.pendingProps), t.child = n, n.return = t;
					e.sibling !== null;
				)
					(e = e.sibling), (n = n.sibling = gi(e, e.pendingProps)), (n.return = t);
				n.sibling = null;
			}
			return t.child;
		}
		function jc(e, t) {
			return (e.lanes & t) === 0 ? ((e = e.dependencies), !!(e !== null && ia(e))) : !0;
		}
		function Mc(e, t, n) {
			switch (t.tag) {
				case 3:
					ve(t, t.stateNode.containerInfo), $i(t, F, e.memoizedState.cache), qi();
					break;
				case 27:
				case 5:
					be(t);
					break;
				case 4:
					ve(t, t.stateNode.containerInfo);
					break;
				case 10:
					$i(t, t.type, t.memoizedProps.value);
					break;
				case 31:
					if (t.memoizedState !== null) return (t.flags |= 128), uo(t), null;
					break;
				case 13:
					var r = t.memoizedState;
					if (r !== null)
						return r.dehydrated === null
							? (n & t.child.childLanes) === 0
								? (lo(t), (e = Ac(e, t, n)), e === null ? null : e.sibling)
								: Cc(e, t, n)
							: (lo(t), (t.flags |= 128), null);
					lo(t);
					break;
				case 19:
					var i = (e.flags & 128) != 0;
					if (
						((r = (n & t.childLanes) !== 0), (r ||= (ra(e, t, n, !1), (n & t.childLanes) !== 0)), i)
					) {
						if (r) return kc(e, t, n);
						t.flags |= 128;
					}
					if (
						((i = t.memoizedState),
						i !== null && ((i.rendering = null), (i.tail = null), (i.lastEffect = null)),
						k(I, I.current),
						r)
					)
						break;
					return null;
				case 22:
					return (t.lanes = 0), lc(e, t, n, t.pendingProps);
				case 24:
					$i(t, F, e.memoizedState.cache);
			}
			return Ac(e, t, n);
		}
		function Nc(e, t, n) {
			if (e !== null)
				if (e.memoizedProps !== t.pendingProps) H = !0;
				else {
					if (!jc(e, n) && !(t.flags & 128)) return (H = !1), Mc(e, t, n);
					H = !!(e.flags & 131072);
				}
			else (H = !1), P && t.flags & 1048576 && Fi(t, Oi, t.index);
			switch (((t.lanes = 0), t.tag)) {
				case 16:
					a: {
						var r = t.pendingProps;
						if (((e = Ma(t.elementType)), (t.type = e), typeof e == `function`))
							hi(e)
								? ((r = Js(e, r)), (t.tag = 1), (t = vc(null, t, e, r, n)))
								: ((t.tag = 0), (t = gc(null, t, e, r, n)));
						else {
							if (e != null) {
								var i = e.$$typeof;
								if (i === C) {
									(t.tag = 11), (t = oc(null, t, e, r, n));
									break a;
								} else if (i === re) {
									(t.tag = 14), (t = sc(null, t, e, r, n));
									break a;
								}
							}
							throw ((t = le(e) || e), Error(s(306, t, ``)));
						}
					}
					return t;
				case 0:
					return gc(e, t, t.type, t.pendingProps, n);
				case 1:
					return (r = t.type), (i = Js(r, t.pendingProps)), vc(e, t, r, i, n);
				case 3:
					a: {
						if ((ve(t, t.stateNode.containerInfo), e === null)) throw Error(s(387));
						r = t.pendingProps;
						var a = t.memoizedState;
						(i = a.element), Ka(e, t), $a(t, r, null, n);
						var o = t.memoizedState;
						if (
							((r = o.cache),
							$i(t, F, r),
							r !== a.cache && na(t, [F], n, !0),
							Qa(),
							(r = o.element),
							a.isDehydrated)
						)
							if (
								((a = { element: r, isDehydrated: !1, cache: o.cache }),
								(t.updateQueue.baseState = a),
								(t.memoizedState = a),
								t.flags & 256)
							) {
								t = yc(e, t, r, n);
								break a;
							} else if (r !== i) {
								(i = wi(Error(s(424)), t)), Yi(i), (t = yc(e, t, r, n));
								break a;
							} else {
								switch (((e = t.stateNode.containerInfo), e.nodeType)) {
									case 9:
										e = e.body;
										break;
									default:
										e = e.nodeName === `HTML` ? e.ownerDocument.body : e;
								}
								for (
									N = cf(e.firstChild),
										zi = t,
										P = !0,
										Bi = null,
										Vi = !0,
										n = Ua(t, null, r, n),
										t.child = n;
									n;
								)
									(n.flags = (n.flags & -3) | 4096), (n = n.sibling);
							}
						else {
							if ((qi(), r === i)) {
								t = Ac(e, t, n);
								break a;
							}
							ac(e, t, r, n);
						}
						t = t.child;
					}
					return t;
				case 26:
					return (
						hc(e, t),
						e === null
							? (n = kf(t.type, null, t.pendingProps, null))
								? (t.memoizedState = n)
								: P ||
									((n = t.type),
									(e = t.pendingProps),
									(r = Bd(ge.current).createElement(n)),
									(r[mt] = t),
									(r[ht] = e),
									Pd(r, n, e),
									Dt(r),
									(t.stateNode = r))
							: (t.memoizedState = kf(t.type, e.memoizedProps, t.pendingProps, e.memoizedState)),
						null
					);
				case 27:
					return (
						be(t),
						e === null &&
							P &&
							((r = t.stateNode = ff(t.type, t.pendingProps, ge.current)),
							(zi = t),
							(Vi = !0),
							(i = N),
							Zd(t.type) ? ((lf = i), (N = cf(r.firstChild))) : (N = i)),
						ac(e, t, t.pendingProps.children, n),
						hc(e, t),
						e === null && (t.flags |= 4194304),
						t.child
					);
				case 5:
					return (
						e === null &&
							P &&
							((i = r = N) &&
								((r = tf(r, t.type, t.pendingProps, Vi)),
								r === null
									? (i = !1)
									: ((t.stateNode = r), (zi = t), (N = cf(r.firstChild)), (Vi = !1), (i = !0))),
							i || Ui(t)),
						be(t),
						(i = t.type),
						(a = t.pendingProps),
						(o = e === null ? null : e.memoizedProps),
						(r = a.children),
						Ud(i, a) ? (r = null) : o !== null && Ud(i, o) && (t.flags |= 32),
						t.memoizedState !== null && ((i = To(e, t, Oo, null, null, n)), (Qf._currentValue = i)),
						hc(e, t),
						ac(e, t, r, n),
						t.child
					);
				case 6:
					return (
						e === null &&
							P &&
							((e = n = N) &&
								((n = nf(n, t.pendingProps, Vi)),
								n === null ? (e = !1) : ((t.stateNode = n), (zi = t), (N = null), (e = !0))),
							e || Ui(t)),
						null
					);
				case 13:
					return Cc(e, t, n);
				case 4:
					return (
						ve(t, t.stateNode.containerInfo),
						(r = t.pendingProps),
						e === null ? (t.child = Ha(t, null, r, n)) : ac(e, t, r, n),
						t.child
					);
				case 11:
					return oc(e, t, t.type, t.pendingProps, n);
				case 7:
					return ac(e, t, t.pendingProps, n), t.child;
				case 8:
					return ac(e, t, t.pendingProps.children, n), t.child;
				case 12:
					return ac(e, t, t.pendingProps.children, n), t.child;
				case 10:
					return (r = t.pendingProps), $i(t, t.type, r.value), ac(e, t, r.children, n), t.child;
				case 9:
					return (
						(i = t.type._context),
						(r = t.pendingProps.children),
						aa(t),
						(i = oa(i)),
						(r = r(i)),
						(t.flags |= 1),
						ac(e, t, r, n),
						t.child
					);
				case 14:
					return sc(e, t, t.type, t.pendingProps, n);
				case 15:
					return cc(e, t, t.type, t.pendingProps, n);
				case 19:
					return kc(e, t, n);
				case 31:
					return mc(e, t, n);
				case 22:
					return lc(e, t, n, t.pendingProps);
				case 24:
					return (
						aa(t),
						(r = oa(F)),
						e === null
							? ((i = Ca()),
								i === null &&
									((i = K),
									(a = fa()),
									(i.pooledCache = a),
									a.refCount++,
									a !== null && (i.pooledCacheLanes |= n),
									(i = a)),
								(t.memoizedState = { parent: r, cache: i }),
								Ga(t),
								$i(t, F, i))
							: ((e.lanes & n) !== 0 && (Ka(e, t), $a(t, null, null, n), Qa()),
								(i = e.memoizedState),
								(a = t.memoizedState),
								i.parent === r
									? ((r = a.cache), $i(t, F, r), r !== i.cache && na(t, [F], n, !0))
									: ((i = { parent: r, cache: r }),
										(t.memoizedState = i),
										t.lanes === 0 && (t.memoizedState = t.updateQueue.baseState = i),
										$i(t, F, r))),
						ac(e, t, t.pendingProps.children, n),
						t.child
					);
				case 29:
					throw t.pendingProps;
			}
			throw Error(s(156, t.tag));
		}
		function Pc(e) {
			e.flags |= 4;
		}
		function Fc(e, t, n, r, i) {
			if (((t = (e.mode & 32) != 0) && (t = !1), t)) {
				if (((e.flags |= 16777216), (i & 335544128) === i))
					if (e.stateNode.complete) e.flags |= 8192;
					else if (wu()) e.flags |= 8192;
					else throw ((Na = ka), Da);
			} else e.flags &= -16777217;
		}
		function Ic(e, t) {
			if (t.type !== `stylesheet` || t.state.loading & 4) e.flags &= -16777217;
			else if (((e.flags |= 16777216), !Wf(t)))
				if (wu()) e.flags |= 8192;
				else throw ((Na = ka), Da);
		}
		function Lc(e, t) {
			t !== null && (e.flags |= 4),
				e.flags & 16384 && ((t = e.tag === 22 ? 536870912 : nt()), (e.lanes |= t), (Yl |= t));
		}
		function Rc(e, t) {
			if (!P)
				switch (e.tailMode) {
					case `hidden`:
						t = e.tail;
						for (var n = null; t !== null; ) t.alternate !== null && (n = t), (t = t.sibling);
						n === null ? (e.tail = null) : (n.sibling = null);
						break;
					case `collapsed`:
						n = e.tail;
						for (var r = null; n !== null; ) n.alternate !== null && (r = n), (n = n.sibling);
						r === null
							? t || e.tail === null
								? (e.tail = null)
								: (e.tail.sibling = null)
							: (r.sibling = null);
				}
		}
		function U(e) {
			var t = e.alternate !== null && e.alternate.child === e.child,
				n = 0,
				r = 0;
			if (t)
				for (var i = e.child; i !== null; )
					(n |= i.lanes | i.childLanes),
						(r |= i.subtreeFlags & 65011712),
						(r |= i.flags & 65011712),
						(i.return = e),
						(i = i.sibling);
			else
				for (i = e.child; i !== null; )
					(n |= i.lanes | i.childLanes),
						(r |= i.subtreeFlags),
						(r |= i.flags),
						(i.return = e),
						(i = i.sibling);
			return (e.subtreeFlags |= r), (e.childLanes = n), t;
		}
		function zc(e, t, n) {
			var r = t.pendingProps;
			switch ((Li(t), t.tag)) {
				case 16:
				case 15:
				case 0:
				case 11:
				case 7:
				case 8:
				case 12:
				case 9:
				case 14:
					return U(t), null;
				case 1:
					return U(t), null;
				case 3:
					return (
						(n = t.stateNode),
						(r = null),
						e !== null && (r = e.memoizedState.cache),
						t.memoizedState.cache !== r && (t.flags |= 2048),
						ea(F),
						ye(),
						n.pendingContext && ((n.context = n.pendingContext), (n.pendingContext = null)),
						(e === null || e.child === null) &&
							(Ki(t)
								? Pc(t)
								: e === null ||
									(e.memoizedState.isDehydrated && !(t.flags & 256)) ||
									((t.flags |= 1024), Ji())),
						U(t),
						null
					);
				case 26:
					var i = t.type,
						a = t.memoizedState;
					return (
						e === null
							? (Pc(t), a === null ? (U(t), Fc(t, i, null, r, n)) : (U(t), Ic(t, a)))
							: a
								? a === e.memoizedState
									? (U(t), (t.flags &= -16777217))
									: (Pc(t), U(t), Ic(t, a))
								: ((e = e.memoizedProps), e !== r && Pc(t), U(t), Fc(t, i, e, r, n)),
						null
					);
				case 27:
					if ((xe(t), (n = ge.current), (i = t.type), e !== null && t.stateNode != null))
						e.memoizedProps !== r && Pc(t);
					else {
						if (!r) {
							if (t.stateNode === null) throw Error(s(166));
							return U(t), null;
						}
						(e = me.current), Ki(t) ? Wi(t, e) : ((e = ff(i, r, n)), (t.stateNode = e), Pc(t));
					}
					return U(t), null;
				case 5:
					if ((xe(t), (i = t.type), e !== null && t.stateNode != null))
						e.memoizedProps !== r && Pc(t);
					else {
						if (!r) {
							if (t.stateNode === null) throw Error(s(166));
							return U(t), null;
						}
						if (((a = me.current), Ki(t))) Wi(t, a);
						else {
							var o = Bd(ge.current);
							switch (a) {
								case 1:
									a = o.createElementNS(`http://www.w3.org/2000/svg`, i);
									break;
								case 2:
									a = o.createElementNS(`http://www.w3.org/1998/Math/MathML`, i);
									break;
								default:
									switch (i) {
										case `svg`:
											a = o.createElementNS(`http://www.w3.org/2000/svg`, i);
											break;
										case `math`:
											a = o.createElementNS(`http://www.w3.org/1998/Math/MathML`, i);
											break;
										case `script`:
											(a = o.createElement(`div`)),
												(a.innerHTML = `<script><\/script>`),
												(a = a.removeChild(a.firstChild));
											break;
										case `select`:
											(a =
												typeof r.is == `string`
													? o.createElement(`select`, { is: r.is })
													: o.createElement(`select`)),
												r.multiple ? (a.multiple = !0) : r.size && (a.size = r.size);
											break;
										default:
											a =
												typeof r.is == `string`
													? o.createElement(i, { is: r.is })
													: o.createElement(i);
									}
							}
							(a[mt] = t), (a[ht] = r);
							a: for (o = t.child; o !== null; ) {
								if (o.tag === 5 || o.tag === 6) a.appendChild(o.stateNode);
								else if (o.tag !== 4 && o.tag !== 27 && o.child !== null) {
									(o.child.return = o), (o = o.child);
									continue;
								}
								if (o === t) break a;
								for (; o.sibling === null; ) {
									if (o.return === null || o.return === t) break a;
									o = o.return;
								}
								(o.sibling.return = o.return), (o = o.sibling);
							}
							t.stateNode = a;
							a: switch ((Pd(a, i, r), i)) {
								case `button`:
								case `input`:
								case `select`:
								case `textarea`:
									r = !!r.autoFocus;
									break a;
								case `img`:
									r = !0;
									break a;
								default:
									r = !1;
							}
							r && Pc(t);
						}
					}
					return U(t), Fc(t, t.type, e === null ? null : e.memoizedProps, t.pendingProps, n), null;
				case 6:
					if (e && t.stateNode != null) e.memoizedProps !== r && Pc(t);
					else {
						if (typeof r != `string` && t.stateNode === null) throw Error(s(166));
						if (((e = ge.current), Ki(t))) {
							if (((e = t.stateNode), (n = t.memoizedProps), (r = null), (i = zi), i !== null))
								switch (i.tag) {
									case 27:
									case 5:
										r = i.memoizedProps;
								}
							(e[mt] = t),
								(e = !!(
									e.nodeValue === n ||
									(r !== null && !0 === r.suppressHydrationWarning) ||
									Md(e.nodeValue, n)
								)),
								e || Ui(t, !0);
						} else (e = Bd(e).createTextNode(r)), (e[mt] = t), (t.stateNode = e);
					}
					return U(t), null;
				case 31:
					if (((n = t.memoizedState), e === null || e.memoizedState !== null)) {
						if (((r = Ki(t)), n !== null)) {
							if (e === null) {
								if (!r) throw Error(s(318));
								if (((e = t.memoizedState), (e = e === null ? null : e.dehydrated), !e))
									throw Error(s(557));
								e[mt] = t;
							} else qi(), !(t.flags & 128) && (t.memoizedState = null), (t.flags |= 4);
							U(t), (e = !1);
						} else
							(n = Ji()),
								e !== null && e.memoizedState !== null && (e.memoizedState.hydrationErrors = n),
								(e = !0);
						if (!e) return t.flags & 256 ? (mo(t), t) : (mo(t), null);
						if (t.flags & 128) throw Error(s(558));
					}
					return U(t), null;
				case 13:
					if (
						((r = t.memoizedState),
						e === null || (e.memoizedState !== null && e.memoizedState.dehydrated !== null))
					) {
						if (((i = Ki(t)), r !== null && r.dehydrated !== null)) {
							if (e === null) {
								if (!i) throw Error(s(318));
								if (((i = t.memoizedState), (i = i === null ? null : i.dehydrated), !i))
									throw Error(s(317));
								i[mt] = t;
							} else qi(), !(t.flags & 128) && (t.memoizedState = null), (t.flags |= 4);
							U(t), (i = !1);
						} else
							(i = Ji()),
								e !== null && e.memoizedState !== null && (e.memoizedState.hydrationErrors = i),
								(i = !0);
						if (!i) return t.flags & 256 ? (mo(t), t) : (mo(t), null);
					}
					return (
						mo(t),
						t.flags & 128
							? ((t.lanes = n), t)
							: ((n = r !== null),
								(e = e !== null && e.memoizedState !== null),
								n &&
									((r = t.child),
									(i = null),
									r.alternate !== null &&
										r.alternate.memoizedState !== null &&
										r.alternate.memoizedState.cachePool !== null &&
										(i = r.alternate.memoizedState.cachePool.pool),
									(a = null),
									r.memoizedState !== null &&
										r.memoizedState.cachePool !== null &&
										(a = r.memoizedState.cachePool.pool),
									a !== i && (r.flags |= 2048)),
								n !== e && n && (t.child.flags |= 8192),
								Lc(t, t.updateQueue),
								U(t),
								null)
					);
				case 4:
					return ye(), e === null && Sd(t.stateNode.containerInfo), U(t), null;
				case 10:
					return ea(t.type), U(t), null;
				case 19:
					if ((O(I), (r = t.memoizedState), r === null)) return U(t), null;
					if (((i = (t.flags & 128) != 0), (a = r.rendering), a === null))
						if (i) Rc(r, !1);
						else {
							if (X !== 0 || (e !== null && e.flags & 128))
								for (e = t.child; e !== null; ) {
									if (((a = ho(e)), a !== null)) {
										for (
											t.flags |= 128,
												Rc(r, !1),
												e = a.updateQueue,
												t.updateQueue = e,
												Lc(t, e),
												t.subtreeFlags = 0,
												e = n,
												n = t.child;
											n !== null;
										)
											_i(n, e), (n = n.sibling);
										return k(I, (I.current & 1) | 2), P && Pi(t, r.treeForkCount), t.child;
									}
									e = e.sibling;
								}
							r.tail !== null &&
								Ne() > tu &&
								((t.flags |= 128), (i = !0), Rc(r, !1), (t.lanes = 4194304));
						}
					else {
						if (!i)
							if (((e = ho(a)), e !== null)) {
								if (
									((t.flags |= 128),
									(i = !0),
									(e = e.updateQueue),
									(t.updateQueue = e),
									Lc(t, e),
									Rc(r, !0),
									r.tail === null && r.tailMode === `hidden` && !a.alternate && !P)
								)
									return U(t), null;
							} else
								2 * Ne() - r.renderingStartTime > tu &&
									n !== 536870912 &&
									((t.flags |= 128), (i = !0), Rc(r, !1), (t.lanes = 4194304));
						r.isBackwards
							? ((a.sibling = t.child), (t.child = a))
							: ((e = r.last), e === null ? (t.child = a) : (e.sibling = a), (r.last = a));
					}
					return r.tail === null
						? (U(t), null)
						: ((e = r.tail),
							(r.rendering = e),
							(r.tail = e.sibling),
							(r.renderingStartTime = Ne()),
							(e.sibling = null),
							(n = I.current),
							k(I, i ? (n & 1) | 2 : n & 1),
							P && Pi(t, r.treeForkCount),
							e);
				case 22:
				case 23:
					return (
						mo(t),
						oo(),
						(r = t.memoizedState !== null),
						e === null
							? r && (t.flags |= 8192)
							: (e.memoizedState !== null) !== r && (t.flags |= 8192),
						r
							? n & 536870912 && !(t.flags & 128) && (U(t), t.subtreeFlags & 6 && (t.flags |= 8192))
							: U(t),
						(n = t.updateQueue),
						n !== null && Lc(t, n.retryQueue),
						(n = null),
						e !== null &&
							e.memoizedState !== null &&
							e.memoizedState.cachePool !== null &&
							(n = e.memoizedState.cachePool.pool),
						(r = null),
						t.memoizedState !== null &&
							t.memoizedState.cachePool !== null &&
							(r = t.memoizedState.cachePool.pool),
						r !== n && (t.flags |= 2048),
						e !== null && O(Sa),
						null
					);
				case 24:
					return (
						(n = null),
						e !== null && (n = e.memoizedState.cache),
						t.memoizedState.cache !== n && (t.flags |= 2048),
						ea(F),
						U(t),
						null
					);
				case 25:
					return null;
				case 30:
					return null;
			}
			throw Error(s(156, t.tag));
		}
		function Bc(e, t) {
			switch ((Li(t), t.tag)) {
				case 1:
					return (e = t.flags), e & 65536 ? ((t.flags = (e & -65537) | 128), t) : null;
				case 3:
					return (
						ea(F),
						ye(),
						(e = t.flags),
						e & 65536 && !(e & 128) ? ((t.flags = (e & -65537) | 128), t) : null
					);
				case 26:
				case 27:
				case 5:
					return xe(t), null;
				case 31:
					if (t.memoizedState !== null) {
						if ((mo(t), t.alternate === null)) throw Error(s(340));
						qi();
					}
					return (e = t.flags), e & 65536 ? ((t.flags = (e & -65537) | 128), t) : null;
				case 13:
					if ((mo(t), (e = t.memoizedState), e !== null && e.dehydrated !== null)) {
						if (t.alternate === null) throw Error(s(340));
						qi();
					}
					return (e = t.flags), e & 65536 ? ((t.flags = (e & -65537) | 128), t) : null;
				case 19:
					return O(I), null;
				case 4:
					return ye(), null;
				case 10:
					return ea(t.type), null;
				case 22:
				case 23:
					return (
						mo(t),
						oo(),
						e !== null && O(Sa),
						(e = t.flags),
						e & 65536 ? ((t.flags = (e & -65537) | 128), t) : null
					);
				case 24:
					return ea(F), null;
				case 25:
					return null;
				default:
					return null;
			}
		}
		function Vc(e, t) {
			switch ((Li(t), t.tag)) {
				case 3:
					ea(F), ye();
					break;
				case 26:
				case 27:
				case 5:
					xe(t);
					break;
				case 4:
					ye();
					break;
				case 31:
					t.memoizedState !== null && mo(t);
					break;
				case 13:
					mo(t);
					break;
				case 19:
					O(I);
					break;
				case 10:
					ea(t.type);
					break;
				case 22:
				case 23:
					mo(t), oo(), e !== null && O(Sa);
					break;
				case 24:
					ea(F);
			}
		}
		function Hc(e, t) {
			try {
				var n = t.updateQueue,
					r = n === null ? null : n.lastEffect;
				if (r !== null) {
					var i = r.next;
					n = i;
					do {
						if ((n.tag & e) === e) {
							r = void 0;
							var a = n.create,
								o = n.inst;
							(r = a()), (o.destroy = r);
						}
						n = n.next;
					} while (n !== i);
				}
			} catch (e) {
				Z(t, t.return, e);
			}
		}
		function Uc(e, t, n) {
			try {
				var r = t.updateQueue,
					i = r === null ? null : r.lastEffect;
				if (i !== null) {
					var a = i.next;
					r = a;
					do {
						if ((r.tag & e) === e) {
							var o = r.inst,
								s = o.destroy;
							if (s !== void 0) {
								(o.destroy = void 0), (i = t);
								var c = n,
									l = s;
								try {
									l();
								} catch (e) {
									Z(i, c, e);
								}
							}
						}
						r = r.next;
					} while (r !== a);
				}
			} catch (e) {
				Z(t, t.return, e);
			}
		}
		function Wc(e) {
			var t = e.updateQueue;
			if (t !== null) {
				var n = e.stateNode;
				try {
					to(t, n);
				} catch (t) {
					Z(e, e.return, t);
				}
			}
		}
		function Gc(e, t, n) {
			(n.props = Js(e.type, e.memoizedProps)), (n.state = e.memoizedState);
			try {
				n.componentWillUnmount();
			} catch (n) {
				Z(e, t, n);
			}
		}
		function Kc(e, t) {
			try {
				var n = e.ref;
				if (n !== null) {
					switch (e.tag) {
						case 26:
						case 27:
						case 5:
							var r = e.stateNode;
							break;
						case 30:
							r = e.stateNode;
							break;
						default:
							r = e.stateNode;
					}
					typeof n == `function` ? (e.refCleanup = n(r)) : (n.current = r);
				}
			} catch (n) {
				Z(e, t, n);
			}
		}
		function qc(e, t) {
			var n = e.ref,
				r = e.refCleanup;
			if (n !== null)
				if (typeof r == `function`)
					try {
						r();
					} catch (n) {
						Z(e, t, n);
					} finally {
						(e.refCleanup = null), (e = e.alternate), e != null && (e.refCleanup = null);
					}
				else if (typeof n == `function`)
					try {
						n(null);
					} catch (n) {
						Z(e, t, n);
					}
				else n.current = null;
		}
		function Jc(e) {
			var t = e.type,
				n = e.memoizedProps,
				r = e.stateNode;
			try {
				a: switch (t) {
					case `button`:
					case `input`:
					case `select`:
					case `textarea`:
						n.autoFocus && r.focus();
						break a;
					case `img`:
						n.src ? (r.src = n.src) : n.srcSet && (r.srcset = n.srcSet);
				}
			} catch (t) {
				Z(e, e.return, t);
			}
		}
		function Yc(e, t, n) {
			try {
				var r = e.stateNode;
				Fd(r, e.type, n, t), (r[ht] = t);
			} catch (t) {
				Z(e, e.return, t);
			}
		}
		function Xc(e) {
			return (
				e.tag === 5 || e.tag === 3 || e.tag === 26 || (e.tag === 27 && Zd(e.type)) || e.tag === 4
			);
		}
		function Zc(e) {
			a: for (;;) {
				for (; e.sibling === null; ) {
					if (e.return === null || Xc(e.return)) return null;
					e = e.return;
				}
				for (
					e.sibling.return = e.return, e = e.sibling;
					e.tag !== 5 && e.tag !== 6 && e.tag !== 18;
				) {
					if ((e.tag === 27 && Zd(e.type)) || e.flags & 2 || e.child === null || e.tag === 4)
						continue a;
					(e.child.return = e), (e = e.child);
				}
				if (!(e.flags & 2)) return e.stateNode;
			}
		}
		function Qc(e, t, n) {
			var r = e.tag;
			if (r === 5 || r === 6)
				(e = e.stateNode),
					t
						? (n.nodeType === 9
								? n.body
								: n.nodeName === `HTML`
									? n.ownerDocument.body
									: n
							).insertBefore(e, t)
						: ((t = n.nodeType === 9 ? n.body : n.nodeName === `HTML` ? n.ownerDocument.body : n),
							t.appendChild(e),
							(n = n._reactRootContainer),
							n != null || t.onclick !== null || (t.onclick = sn));
			else if (
				r !== 4 &&
				(r === 27 && Zd(e.type) && ((n = e.stateNode), (t = null)), (e = e.child), e !== null)
			)
				for (Qc(e, t, n), e = e.sibling; e !== null; ) Qc(e, t, n), (e = e.sibling);
		}
		function $c(e, t, n) {
			var r = e.tag;
			if (r === 5 || r === 6) (e = e.stateNode), t ? n.insertBefore(e, t) : n.appendChild(e);
			else if (r !== 4 && (r === 27 && Zd(e.type) && (n = e.stateNode), (e = e.child), e !== null))
				for ($c(e, t, n), e = e.sibling; e !== null; ) $c(e, t, n), (e = e.sibling);
		}
		function el(e) {
			var t = e.stateNode,
				n = e.memoizedProps;
			try {
				for (var r = e.type, i = t.attributes; i.length; ) t.removeAttributeNode(i[0]);
				Pd(t, r, n), (t[mt] = e), (t[ht] = n);
			} catch (t) {
				Z(e, e.return, t);
			}
		}
		var tl = !1,
			nl = !1,
			rl = !1,
			il = typeof WeakSet == `function` ? WeakSet : Set,
			al = null;
		function ol(e, t) {
			if (((e = e.containerInfo), (Rd = sp), (e = Nr(e)), Pr(e))) {
				if (`selectionStart` in e) var n = { start: e.selectionStart, end: e.selectionEnd };
				else
					a: {
						n = ((n = e.ownerDocument) && n.defaultView) || window;
						var r = n.getSelection && n.getSelection();
						if (r && r.rangeCount !== 0) {
							n = r.anchorNode;
							var i = r.anchorOffset,
								a = r.focusNode;
							r = r.focusOffset;
							try {
								n.nodeType, a.nodeType;
							} catch {
								n = null;
								break a;
							}
							var o = 0,
								c = -1,
								l = -1,
								u = 0,
								d = 0,
								f = e,
								p = null;
							b: for (;;) {
								for (
									var m;
									f !== n || (i !== 0 && f.nodeType !== 3) || (c = o + i),
										f !== a || (r !== 0 && f.nodeType !== 3) || (l = o + r),
										f.nodeType === 3 && (o += f.nodeValue.length),
										(m = f.firstChild) !== null;
								)
									(p = f), (f = m);
								for (;;) {
									if (f === e) break b;
									if (
										(p === n && ++u === i && (c = o),
										p === a && ++d === r && (l = o),
										(m = f.nextSibling) !== null)
									)
										break;
									(f = p), (p = f.parentNode);
								}
								f = m;
							}
							n = c === -1 || l === -1 ? null : { start: c, end: l };
						} else n = null;
					}
				n ||= { start: 0, end: 0 };
			} else n = null;
			for (zd = { focusedElem: e, selectionRange: n }, sp = !1, al = t; al !== null; )
				if (((t = al), (e = t.child), t.subtreeFlags & 1028 && e !== null))
					(e.return = t), (al = e);
				else
					for (; al !== null; ) {
						switch (((t = al), (a = t.alternate), (e = t.flags), t.tag)) {
							case 0:
								if (e & 4 && ((e = t.updateQueue), (e = e === null ? null : e.events), e !== null))
									for (n = 0; n < e.length; n++) (i = e[n]), (i.ref.impl = i.nextImpl);
								break;
							case 11:
							case 15:
								break;
							case 1:
								if (e & 1024 && a !== null) {
									(e = void 0),
										(n = t),
										(i = a.memoizedProps),
										(a = a.memoizedState),
										(r = n.stateNode);
									try {
										var h = Js(n.type, i);
										(e = r.getSnapshotBeforeUpdate(h, a)),
											(r.__reactInternalSnapshotBeforeUpdate = e);
									} catch (e) {
										Z(n, n.return, e);
									}
								}
								break;
							case 3:
								if (e & 1024) {
									if (((e = t.stateNode.containerInfo), (n = e.nodeType), n === 9)) ef(e);
									else if (n === 1)
										switch (e.nodeName) {
											case `HEAD`:
											case `HTML`:
											case `BODY`:
												ef(e);
												break;
											default:
												e.textContent = ``;
										}
								}
								break;
							case 5:
							case 26:
							case 27:
							case 6:
							case 4:
							case 17:
								break;
							default:
								if (e & 1024) throw Error(s(163));
						}
						if (((e = t.sibling), e !== null)) {
							(e.return = t.return), (al = e);
							break;
						}
						al = t.return;
					}
		}
		function sl(e, t, n) {
			var r = n.flags;
			switch (n.tag) {
				case 0:
				case 11:
				case 15:
					xl(e, n), r & 4 && Hc(5, n);
					break;
				case 1:
					if ((xl(e, n), r & 4))
						if (((e = n.stateNode), t === null))
							try {
								e.componentDidMount();
							} catch (e) {
								Z(n, n.return, e);
							}
						else {
							var i = Js(n.type, t.memoizedProps);
							t = t.memoizedState;
							try {
								e.componentDidUpdate(i, t, e.__reactInternalSnapshotBeforeUpdate);
							} catch (e) {
								Z(n, n.return, e);
							}
						}
					r & 64 && Wc(n), r & 512 && Kc(n, n.return);
					break;
				case 3:
					if ((xl(e, n), r & 64 && ((e = n.updateQueue), e !== null))) {
						if (((t = null), n.child !== null))
							switch (n.child.tag) {
								case 27:
								case 5:
									t = n.child.stateNode;
									break;
								case 1:
									t = n.child.stateNode;
							}
						try {
							to(e, t);
						} catch (e) {
							Z(n, n.return, e);
						}
					}
					break;
				case 27:
					t === null && r & 4 && el(n);
				case 26:
				case 5:
					xl(e, n), t === null && r & 4 && Jc(n), r & 512 && Kc(n, n.return);
					break;
				case 12:
					xl(e, n);
					break;
				case 31:
					xl(e, n), r & 4 && fl(e, n);
					break;
				case 13:
					xl(e, n),
						r & 4 && pl(e, n),
						r & 64 &&
							((e = n.memoizedState),
							e !== null && ((e = e.dehydrated), e !== null && ((n = Ju.bind(null, n)), sf(e, n))));
					break;
				case 22:
					if (((r = n.memoizedState !== null || tl), !r)) {
						(t = (t !== null && t.memoizedState !== null) || nl), (i = tl);
						var a = nl;
						(tl = r),
							(nl = t) && !a ? Cl(e, n, (n.subtreeFlags & 8772) != 0) : xl(e, n),
							(tl = i),
							(nl = a);
					}
					break;
				case 30:
					break;
				default:
					xl(e, n);
			}
		}
		function cl(e) {
			var t = e.alternate;
			t !== null && ((e.alternate = null), cl(t)),
				(e.child = null),
				(e.deletions = null),
				(e.sibling = null),
				e.tag === 5 && ((t = e.stateNode), t !== null && St(t)),
				(e.stateNode = null),
				(e.return = null),
				(e.dependencies = null),
				(e.memoizedProps = null),
				(e.memoizedState = null),
				(e.pendingProps = null),
				(e.stateNode = null),
				(e.updateQueue = null);
		}
		var W = null,
			ll = !1;
		function ul(e, t, n) {
			for (n = n.child; n !== null; ) dl(e, t, n), (n = n.sibling);
		}
		function dl(e, t, n) {
			if (Ue && typeof Ue.onCommitFiberUnmount == `function`)
				try {
					Ue.onCommitFiberUnmount(He, n);
				} catch {}
			switch (n.tag) {
				case 26:
					nl || qc(n, t),
						ul(e, t, n),
						n.memoizedState
							? n.memoizedState.count--
							: n.stateNode && ((n = n.stateNode), n.parentNode.removeChild(n));
					break;
				case 27:
					nl || qc(n, t);
					var r = W,
						i = ll;
					Zd(n.type) && ((W = n.stateNode), (ll = !1)),
						ul(e, t, n),
						pf(n.stateNode),
						(W = r),
						(ll = i);
					break;
				case 5:
					nl || qc(n, t);
				case 6:
					if (((r = W), (i = ll), (W = null), ul(e, t, n), (W = r), (ll = i), W !== null))
						if (ll)
							try {
								(W.nodeType === 9
									? W.body
									: W.nodeName === `HTML`
										? W.ownerDocument.body
										: W
								).removeChild(n.stateNode);
							} catch (e) {
								Z(n, t, e);
							}
						else
							try {
								W.removeChild(n.stateNode);
							} catch (e) {
								Z(n, t, e);
							}
					break;
				case 18:
					W !== null &&
						(ll
							? ((e = W),
								Qd(
									e.nodeType === 9 ? e.body : e.nodeName === `HTML` ? e.ownerDocument.body : e,
									n.stateNode,
								),
								Np(e))
							: Qd(W, n.stateNode));
					break;
				case 4:
					(r = W),
						(i = ll),
						(W = n.stateNode.containerInfo),
						(ll = !0),
						ul(e, t, n),
						(W = r),
						(ll = i);
					break;
				case 0:
				case 11:
				case 14:
				case 15:
					Uc(2, n, t), nl || Uc(4, n, t), ul(e, t, n);
					break;
				case 1:
					nl ||
						(qc(n, t),
						(r = n.stateNode),
						typeof r.componentWillUnmount == `function` && Gc(n, t, r)),
						ul(e, t, n);
					break;
				case 21:
					ul(e, t, n);
					break;
				case 22:
					(nl = (r = nl) || n.memoizedState !== null), ul(e, t, n), (nl = r);
					break;
				default:
					ul(e, t, n);
			}
		}
		function fl(e, t) {
			if (
				t.memoizedState === null &&
				((e = t.alternate), e !== null && ((e = e.memoizedState), e !== null))
			) {
				e = e.dehydrated;
				try {
					Np(e);
				} catch (e) {
					Z(t, t.return, e);
				}
			}
		}
		function pl(e, t) {
			if (
				t.memoizedState === null &&
				((e = t.alternate),
				e !== null && ((e = e.memoizedState), e !== null && ((e = e.dehydrated), e !== null)))
			)
				try {
					Np(e);
				} catch (e) {
					Z(t, t.return, e);
				}
		}
		function ml(e) {
			switch (e.tag) {
				case 31:
				case 13:
				case 19:
					var t = e.stateNode;
					return t === null && (t = e.stateNode = new il()), t;
				case 22:
					return (
						(e = e.stateNode), (t = e._retryCache), t === null && (t = e._retryCache = new il()), t
					);
				default:
					throw Error(s(435, e.tag));
			}
		}
		function hl(e, t) {
			var n = ml(e);
			t.forEach(function (t) {
				if (!n.has(t)) {
					n.add(t);
					var r = Yu.bind(null, e, t);
					t.then(r, r);
				}
			});
		}
		function gl(e, t) {
			var n = t.deletions;
			if (n !== null)
				for (var r = 0; r < n.length; r++) {
					var i = n[r],
						a = e,
						o = t,
						c = o;
					a: for (; c !== null; ) {
						switch (c.tag) {
							case 27:
								if (Zd(c.type)) {
									(W = c.stateNode), (ll = !1);
									break a;
								}
								break;
							case 5:
								(W = c.stateNode), (ll = !1);
								break a;
							case 3:
							case 4:
								(W = c.stateNode.containerInfo), (ll = !0);
								break a;
						}
						c = c.return;
					}
					if (W === null) throw Error(s(160));
					dl(a, o, i),
						(W = null),
						(ll = !1),
						(a = i.alternate),
						a !== null && (a.return = null),
						(i.return = null);
				}
			if (t.subtreeFlags & 13886) for (t = t.child; t !== null; ) vl(t, e), (t = t.sibling);
		}
		var _l = null;
		function vl(e, t) {
			var n = e.alternate,
				r = e.flags;
			switch (e.tag) {
				case 0:
				case 11:
				case 14:
				case 15:
					gl(t, e), yl(e), r & 4 && (Uc(3, e, e.return), Hc(3, e), Uc(5, e, e.return));
					break;
				case 1:
					gl(t, e),
						yl(e),
						r & 512 && (nl || n === null || qc(n, n.return)),
						r & 64 &&
							tl &&
							((e = e.updateQueue),
							e !== null &&
								((r = e.callbacks),
								r !== null &&
									((n = e.shared.hiddenCallbacks),
									(e.shared.hiddenCallbacks = n === null ? r : n.concat(r)))));
					break;
				case 26:
					var i = _l;
					if ((gl(t, e), yl(e), r & 512 && (nl || n === null || qc(n, n.return)), r & 4)) {
						var a = n === null ? null : n.memoizedState;
						if (((r = e.memoizedState), n === null))
							if (r === null)
								if (e.stateNode === null) {
									a: {
										(r = e.type), (n = e.memoizedProps), (i = i.ownerDocument || i);
										b: switch (r) {
											case `title`:
												(a = i.getElementsByTagName(`title`)[0]),
													(!a ||
														a[xt] ||
														a[mt] ||
														a.namespaceURI === `http://www.w3.org/2000/svg` ||
														a.hasAttribute(`itemprop`)) &&
														((a = i.createElement(r)),
														i.head.insertBefore(a, i.querySelector(`head > title`))),
													Pd(a, r, n),
													(a[mt] = e),
													Dt(a),
													(r = a);
												break a;
											case `link`:
												var o = Vf(`link`, `href`, i).get(r + (n.href || ``));
												if (o) {
													for (var c = 0; c < o.length; c++)
														if (
															((a = o[c]),
															a.getAttribute(`href`) ===
																(n.href == null || n.href === `` ? null : n.href) &&
																a.getAttribute(`rel`) === (n.rel == null ? null : n.rel) &&
																a.getAttribute(`title`) === (n.title == null ? null : n.title) &&
																a.getAttribute(`crossorigin`) ===
																	(n.crossOrigin == null ? null : n.crossOrigin))
														) {
															o.splice(c, 1);
															break b;
														}
												}
												(a = i.createElement(r)), Pd(a, r, n), i.head.appendChild(a);
												break;
											case `meta`:
												if ((o = Vf(`meta`, `content`, i).get(r + (n.content || ``)))) {
													for (c = 0; c < o.length; c++)
														if (
															((a = o[c]),
															a.getAttribute(`content`) ===
																(n.content == null ? null : `` + n.content) &&
																a.getAttribute(`name`) === (n.name == null ? null : n.name) &&
																a.getAttribute(`property`) ===
																	(n.property == null ? null : n.property) &&
																a.getAttribute(`http-equiv`) ===
																	(n.httpEquiv == null ? null : n.httpEquiv) &&
																a.getAttribute(`charset`) ===
																	(n.charSet == null ? null : n.charSet))
														) {
															o.splice(c, 1);
															break b;
														}
												}
												(a = i.createElement(r)), Pd(a, r, n), i.head.appendChild(a);
												break;
											default:
												throw Error(s(468, r));
										}
										(a[mt] = e), Dt(a), (r = a);
									}
									e.stateNode = r;
								} else Hf(i, e.type, e.stateNode);
							else e.stateNode = If(i, r, e.memoizedProps);
						else
							a === r
								? r === null && e.stateNode !== null && Yc(e, e.memoizedProps, n.memoizedProps)
								: (a === null
										? n.stateNode !== null && ((n = n.stateNode), n.parentNode.removeChild(n))
										: a.count--,
									r === null ? Hf(i, e.type, e.stateNode) : If(i, r, e.memoizedProps));
					}
					break;
				case 27:
					gl(t, e),
						yl(e),
						r & 512 && (nl || n === null || qc(n, n.return)),
						n !== null && r & 4 && Yc(e, e.memoizedProps, n.memoizedProps);
					break;
				case 5:
					if ((gl(t, e), yl(e), r & 512 && (nl || n === null || qc(n, n.return)), e.flags & 32)) {
						i = e.stateNode;
						try {
							j(i, ``);
						} catch (t) {
							Z(e, e.return, t);
						}
					}
					r & 4 &&
						e.stateNode != null &&
						((i = e.memoizedProps), Yc(e, i, n === null ? i : n.memoizedProps)),
						r & 1024 && (rl = !0);
					break;
				case 6:
					if ((gl(t, e), yl(e), r & 4)) {
						if (e.stateNode === null) throw Error(s(162));
						(r = e.memoizedProps), (n = e.stateNode);
						try {
							n.nodeValue = r;
						} catch (t) {
							Z(e, e.return, t);
						}
					}
					break;
				case 3:
					if (
						((Bf = null),
						(i = _l),
						(_l = gf(t.containerInfo)),
						gl(t, e),
						(_l = i),
						yl(e),
						r & 4 && n !== null && n.memoizedState.isDehydrated)
					)
						try {
							Np(t.containerInfo);
						} catch (t) {
							Z(e, e.return, t);
						}
					rl && ((rl = !1), bl(e));
					break;
				case 4:
					(r = _l), (_l = gf(e.stateNode.containerInfo)), gl(t, e), yl(e), (_l = r);
					break;
				case 12:
					gl(t, e), yl(e);
					break;
				case 31:
					gl(t, e),
						yl(e),
						r & 4 && ((r = e.updateQueue), r !== null && ((e.updateQueue = null), hl(e, r)));
					break;
				case 13:
					gl(t, e),
						yl(e),
						e.child.flags & 8192 &&
							(e.memoizedState !== null) != (n !== null && n.memoizedState !== null) &&
							($l = Ne()),
						r & 4 && ((r = e.updateQueue), r !== null && ((e.updateQueue = null), hl(e, r)));
					break;
				case 22:
					i = e.memoizedState !== null;
					var l = n !== null && n.memoizedState !== null,
						u = tl,
						d = nl;
					if (((tl = u || i), (nl = d || l), gl(t, e), (nl = d), (tl = u), yl(e), r & 8192))
						a: for (
							t = e.stateNode,
								t._visibility = i ? t._visibility & -2 : t._visibility | 1,
								i && (n === null || l || tl || nl || Sl(e)),
								n = null,
								t = e;
							;
						) {
							if (t.tag === 5 || t.tag === 26) {
								if (n === null) {
									l = n = t;
									try {
										if (((a = l.stateNode), i))
											(o = a.style),
												typeof o.setProperty == `function`
													? o.setProperty(`display`, `none`, `important`)
													: (o.display = `none`);
										else {
											c = l.stateNode;
											var f = l.memoizedProps.style,
												p = f != null && f.hasOwnProperty(`display`) ? f.display : null;
											c.style.display = p == null || typeof p == `boolean` ? `` : (`` + p).trim();
										}
									} catch (e) {
										Z(l, l.return, e);
									}
								}
							} else if (t.tag === 6) {
								if (n === null) {
									l = t;
									try {
										l.stateNode.nodeValue = i ? `` : l.memoizedProps;
									} catch (e) {
										Z(l, l.return, e);
									}
								}
							} else if (t.tag === 18) {
								if (n === null) {
									l = t;
									try {
										var m = l.stateNode;
										i ? $d(m, !0) : $d(l.stateNode, !1);
									} catch (e) {
										Z(l, l.return, e);
									}
								}
							} else if (
								((t.tag !== 22 && t.tag !== 23) || t.memoizedState === null || t === e) &&
								t.child !== null
							) {
								(t.child.return = t), (t = t.child);
								continue;
							}
							if (t === e) break a;
							for (; t.sibling === null; ) {
								if (t.return === null || t.return === e) break a;
								n === t && (n = null), (t = t.return);
							}
							n === t && (n = null), (t.sibling.return = t.return), (t = t.sibling);
						}
					r & 4 &&
						((r = e.updateQueue),
						r !== null && ((n = r.retryQueue), n !== null && ((r.retryQueue = null), hl(e, n))));
					break;
				case 19:
					gl(t, e),
						yl(e),
						r & 4 && ((r = e.updateQueue), r !== null && ((e.updateQueue = null), hl(e, r)));
					break;
				case 30:
					break;
				case 21:
					break;
				default:
					gl(t, e), yl(e);
			}
		}
		function yl(e) {
			var t = e.flags;
			if (t & 2) {
				try {
					for (var n, r = e.return; r !== null; ) {
						if (Xc(r)) {
							n = r;
							break;
						}
						r = r.return;
					}
					if (n == null) throw Error(s(160));
					switch (n.tag) {
						case 27:
							var i = n.stateNode;
							$c(e, Zc(e), i);
							break;
						case 5:
							var a = n.stateNode;
							n.flags & 32 && (j(a, ``), (n.flags &= -33)), $c(e, Zc(e), a);
							break;
						case 3:
						case 4:
							var o = n.stateNode.containerInfo;
							Qc(e, Zc(e), o);
							break;
						default:
							throw Error(s(161));
					}
				} catch (t) {
					Z(e, e.return, t);
				}
				e.flags &= -3;
			}
			t & 4096 && (e.flags &= -4097);
		}
		function bl(e) {
			if (e.subtreeFlags & 1024)
				for (e = e.child; e !== null; ) {
					var t = e;
					bl(t), t.tag === 5 && t.flags & 1024 && t.stateNode.reset(), (e = e.sibling);
				}
		}
		function xl(e, t) {
			if (t.subtreeFlags & 8772)
				for (t = t.child; t !== null; ) sl(e, t.alternate, t), (t = t.sibling);
		}
		function Sl(e) {
			for (e = e.child; e !== null; ) {
				var t = e;
				switch (t.tag) {
					case 0:
					case 11:
					case 14:
					case 15:
						Uc(4, t, t.return), Sl(t);
						break;
					case 1:
						qc(t, t.return);
						var n = t.stateNode;
						typeof n.componentWillUnmount == `function` && Gc(t, t.return, n), Sl(t);
						break;
					case 27:
						pf(t.stateNode);
					case 26:
					case 5:
						qc(t, t.return), Sl(t);
						break;
					case 22:
						t.memoizedState === null && Sl(t);
						break;
					case 30:
						Sl(t);
						break;
					default:
						Sl(t);
				}
				e = e.sibling;
			}
		}
		function Cl(e, t, n) {
			for (n &&= (t.subtreeFlags & 8772) != 0, t = t.child; t !== null; ) {
				var r = t.alternate,
					i = e,
					a = t,
					o = a.flags;
				switch (a.tag) {
					case 0:
					case 11:
					case 15:
						Cl(i, a, n), Hc(4, a);
						break;
					case 1:
						if ((Cl(i, a, n), (r = a), (i = r.stateNode), typeof i.componentDidMount == `function`))
							try {
								i.componentDidMount();
							} catch (e) {
								Z(r, r.return, e);
							}
						if (((r = a), (i = r.updateQueue), i !== null)) {
							var s = r.stateNode;
							try {
								var c = i.shared.hiddenCallbacks;
								if (c !== null)
									for (i.shared.hiddenCallbacks = null, i = 0; i < c.length; i++) eo(c[i], s);
							} catch (e) {
								Z(r, r.return, e);
							}
						}
						n && o & 64 && Wc(a), Kc(a, a.return);
						break;
					case 27:
						el(a);
					case 26:
					case 5:
						Cl(i, a, n), n && r === null && o & 4 && Jc(a), Kc(a, a.return);
						break;
					case 12:
						Cl(i, a, n);
						break;
					case 31:
						Cl(i, a, n), n && o & 4 && fl(i, a);
						break;
					case 13:
						Cl(i, a, n), n && o & 4 && pl(i, a);
						break;
					case 22:
						a.memoizedState === null && Cl(i, a, n), Kc(a, a.return);
						break;
					case 30:
						break;
					default:
						Cl(i, a, n);
				}
				t = t.sibling;
			}
		}
		function wl(e, t) {
			var n = null;
			e !== null &&
				e.memoizedState !== null &&
				e.memoizedState.cachePool !== null &&
				(n = e.memoizedState.cachePool.pool),
				(e = null),
				t.memoizedState !== null &&
					t.memoizedState.cachePool !== null &&
					(e = t.memoizedState.cachePool.pool),
				e !== n && (e != null && e.refCount++, n != null && pa(n));
		}
		function Tl(e, t) {
			(e = null),
				t.alternate !== null && (e = t.alternate.memoizedState.cache),
				(t = t.memoizedState.cache),
				t !== e && (t.refCount++, e != null && pa(e));
		}
		function El(e, t, n, r) {
			if (t.subtreeFlags & 10256) for (t = t.child; t !== null; ) Dl(e, t, n, r), (t = t.sibling);
		}
		function Dl(e, t, n, r) {
			var i = t.flags;
			switch (t.tag) {
				case 0:
				case 11:
				case 15:
					El(e, t, n, r), i & 2048 && Hc(9, t);
					break;
				case 1:
					El(e, t, n, r);
					break;
				case 3:
					El(e, t, n, r),
						i & 2048 &&
							((e = null),
							t.alternate !== null && (e = t.alternate.memoizedState.cache),
							(t = t.memoizedState.cache),
							t !== e && (t.refCount++, e != null && pa(e)));
					break;
				case 12:
					if (i & 2048) {
						El(e, t, n, r), (e = t.stateNode);
						try {
							var a = t.memoizedProps,
								o = a.id,
								s = a.onPostCommit;
							typeof s == `function` &&
								s(o, t.alternate === null ? `mount` : `update`, e.passiveEffectDuration, -0);
						} catch (e) {
							Z(t, t.return, e);
						}
					} else El(e, t, n, r);
					break;
				case 31:
					El(e, t, n, r);
					break;
				case 13:
					El(e, t, n, r);
					break;
				case 23:
					break;
				case 22:
					(a = t.stateNode),
						(o = t.alternate),
						t.memoizedState === null
							? a._visibility & 2
								? El(e, t, n, r)
								: ((a._visibility |= 2), Ol(e, t, n, r, (t.subtreeFlags & 10256) != 0 || !1))
							: a._visibility & 2
								? El(e, t, n, r)
								: kl(e, t),
						i & 2048 && wl(o, t);
					break;
				case 24:
					El(e, t, n, r), i & 2048 && Tl(t.alternate, t);
					break;
				default:
					El(e, t, n, r);
			}
		}
		function Ol(e, t, n, r, i) {
			for (i &&= (t.subtreeFlags & 10256) != 0 || !1, t = t.child; t !== null; ) {
				var a = e,
					o = t,
					s = n,
					c = r,
					l = o.flags;
				switch (o.tag) {
					case 0:
					case 11:
					case 15:
						Ol(a, o, s, c, i), Hc(8, o);
						break;
					case 23:
						break;
					case 22:
						var u = o.stateNode;
						o.memoizedState === null
							? ((u._visibility |= 2), Ol(a, o, s, c, i))
							: u._visibility & 2
								? Ol(a, o, s, c, i)
								: kl(a, o),
							i && l & 2048 && wl(o.alternate, o);
						break;
					case 24:
						Ol(a, o, s, c, i), i && l & 2048 && Tl(o.alternate, o);
						break;
					default:
						Ol(a, o, s, c, i);
				}
				t = t.sibling;
			}
		}
		function kl(e, t) {
			if (t.subtreeFlags & 10256)
				for (t = t.child; t !== null; ) {
					var n = e,
						r = t,
						i = r.flags;
					switch (r.tag) {
						case 22:
							kl(n, r), i & 2048 && wl(r.alternate, r);
							break;
						case 24:
							kl(n, r), i & 2048 && Tl(r.alternate, r);
							break;
						default:
							kl(n, r);
					}
					t = t.sibling;
				}
		}
		var Al = 8192;
		function jl(e, t, n) {
			if (e.subtreeFlags & Al) for (e = e.child; e !== null; ) Ml(e, t, n), (e = e.sibling);
		}
		function Ml(e, t, n) {
			switch (e.tag) {
				case 26:
					jl(e, t, n),
						e.flags & Al && e.memoizedState !== null && Gf(n, _l, e.memoizedState, e.memoizedProps);
					break;
				case 5:
					jl(e, t, n);
					break;
				case 3:
				case 4:
					var r = _l;
					(_l = gf(e.stateNode.containerInfo)), jl(e, t, n), (_l = r);
					break;
				case 22:
					e.memoizedState === null &&
						((r = e.alternate),
						r !== null && r.memoizedState !== null
							? ((r = Al), (Al = 16777216), jl(e, t, n), (Al = r))
							: jl(e, t, n));
					break;
				default:
					jl(e, t, n);
			}
		}
		function Nl(e) {
			var t = e.alternate;
			if (t !== null && ((e = t.child), e !== null)) {
				t.child = null;
				do (t = e.sibling), (e.sibling = null), (e = t);
				while (e !== null);
			}
		}
		function Pl(e) {
			var t = e.deletions;
			if (e.flags & 16) {
				if (t !== null)
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						(al = r), Ll(r, e);
					}
				Nl(e);
			}
			if (e.subtreeFlags & 10256) for (e = e.child; e !== null; ) Fl(e), (e = e.sibling);
		}
		function Fl(e) {
			switch (e.tag) {
				case 0:
				case 11:
				case 15:
					Pl(e), e.flags & 2048 && Uc(9, e, e.return);
					break;
				case 3:
					Pl(e);
					break;
				case 12:
					Pl(e);
					break;
				case 22:
					var t = e.stateNode;
					e.memoizedState !== null &&
					t._visibility & 2 &&
					(e.return === null || e.return.tag !== 13)
						? ((t._visibility &= -3), Il(e))
						: Pl(e);
					break;
				default:
					Pl(e);
			}
		}
		function Il(e) {
			var t = e.deletions;
			if (e.flags & 16) {
				if (t !== null)
					for (var n = 0; n < t.length; n++) {
						var r = t[n];
						(al = r), Ll(r, e);
					}
				Nl(e);
			}
			for (e = e.child; e !== null; ) {
				switch (((t = e), t.tag)) {
					case 0:
					case 11:
					case 15:
						Uc(8, t, t.return), Il(t);
						break;
					case 22:
						(n = t.stateNode), n._visibility & 2 && ((n._visibility &= -3), Il(t));
						break;
					default:
						Il(t);
				}
				e = e.sibling;
			}
		}
		function Ll(e, t) {
			for (; al !== null; ) {
				var n = al;
				switch (n.tag) {
					case 0:
					case 11:
					case 15:
						Uc(8, n, t);
						break;
					case 23:
					case 22:
						if (n.memoizedState !== null && n.memoizedState.cachePool !== null) {
							var r = n.memoizedState.cachePool.pool;
							r != null && r.refCount++;
						}
						break;
					case 24:
						pa(n.memoizedState.cache);
				}
				if (((r = n.child), r !== null)) (r.return = n), (al = r);
				else
					a: for (n = e; al !== null; ) {
						r = al;
						var i = r.sibling,
							a = r.return;
						if ((cl(r), r === n)) {
							al = null;
							break a;
						}
						if (i !== null) {
							(i.return = a), (al = i);
							break a;
						}
						al = a;
					}
			}
		}
		var Rl = {
				getCacheForType: function (e) {
					var t = oa(F),
						n = t.data.get(e);
					return n === void 0 && ((n = e()), t.data.set(e, n)), n;
				},
				cacheSignal: function () {
					return oa(F).controller.signal;
				},
			},
			zl = typeof WeakMap == `function` ? WeakMap : Map,
			G = 0,
			K = null,
			q = null,
			J = 0,
			Y = 0,
			Bl = null,
			Vl = !1,
			Hl = !1,
			Ul = !1,
			Wl = 0,
			X = 0,
			Gl = 0,
			Kl = 0,
			ql = 0,
			Jl = 0,
			Yl = 0,
			Xl = null,
			Zl = null,
			Ql = !1,
			$l = 0,
			eu = 0,
			tu = 1 / 0,
			nu = null,
			ru = null,
			iu = 0,
			au = null,
			ou = null,
			su = 0,
			cu = 0,
			lu = null,
			uu = null,
			du = 0,
			fu = null;
		function pu() {
			return G & 2 && J !== 0 ? J & -J : T.T === null ? dt() : dd();
		}
		function mu() {
			if (Jl === 0)
				if (!(J & 536870912) || P) {
					var e = Xe;
					(Xe <<= 1), !(Xe & 3932160) && (Xe = 262144), (Jl = e);
				} else Jl = 536870912;
			return (e = so.current), e !== null && (e.flags |= 32), Jl;
		}
		function hu(e, t, n) {
			((e === K && (Y === 2 || Y === 9)) || e.cancelPendingCommit !== null) &&
				(Su(e, 0), yu(e, J, Jl, !1)),
				it(e, n),
				(!(G & 2) || e !== K) &&
					(e === K && (!(G & 2) && (Kl |= n), X === 4 && yu(e, J, Jl, !1)), rd(e));
		}
		function gu(e, t, n) {
			if (G & 6) throw Error(s(327));
			var r = (!n && (t & 127) == 0 && (t & e.expiredLanes) === 0) || et(e, t),
				i = r ? Au(e, t) : Ou(e, t, !0),
				a = r;
			do {
				if (i === 0) {
					Hl && !r && yu(e, t, 0, !1);
					break;
				} else {
					if (((n = e.current.alternate), a && !vu(n))) {
						(i = Ou(e, t, !1)), (a = !1);
						continue;
					}
					if (i === 2) {
						if (((a = t), e.errorRecoveryDisabledLanes & a)) var o = 0;
						else
							(o = e.pendingLanes & -536870913),
								(o = o === 0 ? (o & 536870912 ? 536870912 : 0) : o);
						if (o !== 0) {
							t = o;
							a: {
								var c = e;
								i = Xl;
								var l = c.current.memoizedState.isDehydrated;
								if ((l && (Su(c, o).flags |= 256), (o = Ou(c, o, !1)), o !== 2)) {
									if (Ul && !l) {
										(c.errorRecoveryDisabledLanes |= a), (Kl |= a), (i = 4);
										break a;
									}
									(a = Zl), (Zl = i), a !== null && (Zl === null ? (Zl = a) : Zl.push.apply(Zl, a));
								}
								i = o;
							}
							if (((a = !1), i !== 2)) continue;
						}
					}
					if (i === 1) {
						Su(e, 0), yu(e, t, 0, !0);
						break;
					}
					a: {
						switch (((r = e), (a = i), a)) {
							case 0:
							case 1:
								throw Error(s(345));
							case 4:
								if ((t & 4194048) !== t) break;
							case 6:
								yu(r, t, Jl, !Vl);
								break a;
							case 2:
								Zl = null;
								break;
							case 3:
							case 5:
								break;
							default:
								throw Error(s(329));
						}
						if ((t & 62914560) === t && ((i = $l + 300 - Ne()), 10 < i)) {
							if ((yu(r, t, Jl, !Vl), $e(r, 0, !0) !== 0)) break a;
							(su = t),
								(r.timeoutHandle = Kd(
									_u.bind(null, r, n, Zl, nu, Ql, t, Jl, Kl, Yl, Vl, a, `Throttled`, -0, 0),
									i,
								));
							break a;
						}
						_u(r, n, Zl, nu, Ql, t, Jl, Kl, Yl, Vl, a, null, -0, 0);
					}
				}
				break;
			} while (1);
			rd(e);
		}
		function _u(e, t, n, r, i, a, o, s, c, l, u, d, f, p) {
			if (((e.timeoutHandle = -1), (d = t.subtreeFlags), d & 8192 || (d & 16785408) == 16785408)) {
				(d = {
					stylesheets: null,
					count: 0,
					imgCount: 0,
					imgBytes: 0,
					suspenseyImages: [],
					waitingForImages: !0,
					waitingForViewTransition: !1,
					unsuspend: sn,
				}),
					Ml(t, a, d);
				var m = (a & 62914560) === a ? $l - Ne() : (a & 4194048) === a ? eu - Ne() : 0;
				if (((m = qf(d, m)), m !== null)) {
					(su = a),
						(e.cancelPendingCommit = m(Lu.bind(null, e, t, a, n, r, i, o, s, c, u, d, null, f, p))),
						yu(e, a, o, !l);
					return;
				}
			}
			Lu(e, t, a, n, r, i, o, s, c);
		}
		function vu(e) {
			for (var t = e; ; ) {
				var n = t.tag;
				if (
					(n === 0 || n === 11 || n === 15) &&
					t.flags & 16384 &&
					((n = t.updateQueue), n !== null && ((n = n.stores), n !== null))
				)
					for (var r = 0; r < n.length; r++) {
						var i = n[r],
							a = i.getSnapshot;
						i = i.value;
						try {
							if (!Or(a(), i)) return !1;
						} catch {
							return !1;
						}
					}
				if (((n = t.child), t.subtreeFlags & 16384 && n !== null)) (n.return = t), (t = n);
				else {
					if (t === e) break;
					for (; t.sibling === null; ) {
						if (t.return === null || t.return === e) return !0;
						t = t.return;
					}
					(t.sibling.return = t.return), (t = t.sibling);
				}
			}
			return !0;
		}
		function yu(e, t, n, r) {
			(t &= ~ql),
				(t &= ~Kl),
				(e.suspendedLanes |= t),
				(e.pingedLanes &= ~t),
				r && (e.warmLanes |= t),
				(r = e.expirationTimes);
			for (var i = t; 0 < i; ) {
				var a = 31 - Ge(i),
					o = 1 << a;
				(r[a] = -1), (i &= ~o);
			}
			n !== 0 && ot(e, n, t);
		}
		function bu() {
			return G & 6 ? !0 : (id(0, !1), !1);
		}
		function xu() {
			if (q !== null) {
				if (Y === 0) var e = q.return;
				else (e = q), (Qi = Zi = null), jo(e), (Ia = null), (La = 0), (e = q);
				for (; e !== null; ) Vc(e.alternate, e), (e = e.return);
				q = null;
			}
		}
		function Su(e, t) {
			var n = e.timeoutHandle;
			n !== -1 && ((e.timeoutHandle = -1), qd(n)),
				(n = e.cancelPendingCommit),
				n !== null && ((e.cancelPendingCommit = null), n()),
				(su = 0),
				xu(),
				(K = e),
				(q = n = gi(e.current, null)),
				(J = t),
				(Y = 0),
				(Bl = null),
				(Vl = !1),
				(Hl = et(e, t)),
				(Ul = !1),
				(Yl = Jl = ql = Kl = Gl = X = 0),
				(Zl = Xl = null),
				(Ql = !1),
				t & 8 && (t |= t & 32);
			var r = e.entangledLanes;
			if (r !== 0)
				for (e = e.entanglements, r &= t; 0 < r; ) {
					var i = 31 - Ge(r),
						a = 1 << i;
					(t |= e[i]), (r &= ~a);
				}
			return (Wl = t), oi(), n;
		}
		function Cu(e, t) {
			(L = null),
				(T.H = Bs),
				t === Ea || t === Oa
					? ((t = Pa()), (Y = 3))
					: t === Da
						? ((t = Pa()), (Y = 4))
						: (Y = t === ic ? 8 : typeof t == `object` && t && typeof t.then == `function` ? 6 : 1),
				(Bl = t),
				q === null && ((X = 1), Qs(e, wi(t, e.current)));
		}
		function wu() {
			var e = so.current;
			return e === null
				? !0
				: (J & 4194048) === J
					? co === null
					: (J & 62914560) === J || J & 536870912
						? e === co
						: !1;
		}
		function Tu() {
			var e = T.H;
			return (T.H = Bs), e === null ? Bs : e;
		}
		function Eu() {
			var e = T.A;
			return (T.A = Rl), e;
		}
		function Du() {
			(X = 4),
				Vl || ((J & 4194048) !== J && so.current !== null) || (Hl = !0),
				(!(Gl & 134217727) && !(Kl & 134217727)) || K === null || yu(K, J, Jl, !1);
		}
		function Ou(e, t, n) {
			var r = G;
			G |= 2;
			var i = Tu(),
				a = Eu();
			(K !== e || J !== t) && ((nu = null), Su(e, t)), (t = !1);
			var o = X;
			a: do
				try {
					if (Y !== 0 && q !== null) {
						var s = q,
							c = Bl;
						switch (Y) {
							case 8:
								xu(), (o = 6);
								break a;
							case 3:
							case 2:
							case 9:
							case 6:
								so.current === null && (t = !0);
								var l = Y;
								if (((Y = 0), (Bl = null), Pu(e, s, c, l), n && Hl)) {
									o = 0;
									break a;
								}
								break;
							default:
								(l = Y), (Y = 0), (Bl = null), Pu(e, s, c, l);
						}
					}
					ku(), (o = X);
					break;
				} catch (t) {
					Cu(e, t);
				}
			while (1);
			return (
				t && e.shellSuspendCounter++,
				(Qi = Zi = null),
				(G = r),
				(T.H = i),
				(T.A = a),
				q === null && ((K = null), (J = 0), oi()),
				o
			);
		}
		function ku() {
			for (; q !== null; ) Mu(q);
		}
		function Au(e, t) {
			var n = G;
			G |= 2;
			var r = Tu(),
				i = Eu();
			K !== e || J !== t ? ((nu = null), (tu = Ne() + 500), Su(e, t)) : (Hl = et(e, t));
			a: do
				try {
					if (Y !== 0 && q !== null) {
						t = q;
						var a = Bl;
						b: switch (Y) {
							case 1:
								(Y = 0), (Bl = null), Pu(e, t, a, 1);
								break;
							case 2:
							case 9:
								if (Aa(a)) {
									(Y = 0), (Bl = null), Nu(t);
									break;
								}
								(t = function () {
									(Y !== 2 && Y !== 9) || K !== e || (Y = 7), rd(e);
								}),
									a.then(t, t);
								break a;
							case 3:
								Y = 7;
								break a;
							case 4:
								Y = 5;
								break a;
							case 7:
								Aa(a) ? ((Y = 0), (Bl = null), Nu(t)) : ((Y = 0), (Bl = null), Pu(e, t, a, 7));
								break;
							case 5:
								var o = null;
								switch (q.tag) {
									case 26:
										o = q.memoizedState;
									case 5:
									case 27:
										var c = q;
										if (o ? Wf(o) : c.stateNode.complete) {
											(Y = 0), (Bl = null);
											var l = c.sibling;
											if (l !== null) q = l;
											else {
												var u = c.return;
												u === null ? (q = null) : ((q = u), Fu(u));
											}
											break b;
										}
								}
								(Y = 0), (Bl = null), Pu(e, t, a, 5);
								break;
							case 6:
								(Y = 0), (Bl = null), Pu(e, t, a, 6);
								break;
							case 8:
								xu(), (X = 6);
								break a;
							default:
								throw Error(s(462));
						}
					}
					ju();
					break;
				} catch (t) {
					Cu(e, t);
				}
			while (1);
			return (
				(Qi = Zi = null),
				(T.H = r),
				(T.A = i),
				(G = n),
				q === null ? ((K = null), (J = 0), oi(), X) : 0
			);
		}
		function ju() {
			for (; q !== null && !je(); ) Mu(q);
		}
		function Mu(e) {
			var t = Nc(e.alternate, e, Wl);
			(e.memoizedProps = e.pendingProps), t === null ? Fu(e) : (q = t);
		}
		function Nu(e) {
			var t = e,
				n = t.alternate;
			switch (t.tag) {
				case 15:
				case 0:
					t = _c(n, t, t.pendingProps, t.type, void 0, J);
					break;
				case 11:
					t = _c(n, t, t.pendingProps, t.type.render, t.ref, J);
					break;
				case 5:
					jo(t);
				default:
					Vc(n, t), (t = q = _i(t, Wl)), (t = Nc(n, t, Wl));
			}
			(e.memoizedProps = e.pendingProps), t === null ? Fu(e) : (q = t);
		}
		function Pu(e, t, n, r) {
			(Qi = Zi = null), jo(t), (Ia = null), (La = 0);
			var i = t.return;
			try {
				if (rc(e, i, t, n, J)) {
					(X = 1), Qs(e, wi(n, e.current)), (q = null);
					return;
				}
			} catch (t) {
				if (i !== null) throw ((q = i), t);
				(X = 1), Qs(e, wi(n, e.current)), (q = null);
				return;
			}
			t.flags & 32768
				? (P || r === 1
						? (e = !0)
						: Hl || J & 536870912
							? (e = !1)
							: ((Vl = e = !0),
								(r === 2 || r === 9 || r === 3 || r === 6) &&
									((r = so.current), r !== null && r.tag === 13 && (r.flags |= 16384))),
					Iu(t, e))
				: Fu(t);
		}
		function Fu(e) {
			var t = e;
			do {
				if (t.flags & 32768) {
					Iu(t, Vl);
					return;
				}
				e = t.return;
				var n = zc(t.alternate, t, Wl);
				if (n !== null) {
					q = n;
					return;
				}
				if (((t = t.sibling), t !== null)) {
					q = t;
					return;
				}
				q = t = e;
			} while (t !== null);
			X === 0 && (X = 5);
		}
		function Iu(e, t) {
			do {
				var n = Bc(e.alternate, e);
				if (n !== null) {
					(n.flags &= 32767), (q = n);
					return;
				}
				if (
					((n = e.return),
					n !== null && ((n.flags |= 32768), (n.subtreeFlags = 0), (n.deletions = null)),
					!t && ((e = e.sibling), e !== null))
				) {
					q = e;
					return;
				}
				q = e = n;
			} while (e !== null);
			(X = 6), (q = null);
		}
		function Lu(e, t, n, r, i, a, o, c, l) {
			e.cancelPendingCommit = null;
			do Hu();
			while (iu !== 0);
			if (G & 6) throw Error(s(327));
			if (t !== null) {
				if (t === e.current) throw Error(s(177));
				if (
					((a = t.lanes | t.childLanes),
					(a |= ai),
					at(e, n, a, o, c, l),
					e === K && ((q = K = null), (J = 0)),
					(ou = t),
					(au = e),
					(su = n),
					(cu = a),
					(lu = i),
					(uu = r),
					t.subtreeFlags & 10256 || t.flags & 10256
						? ((e.callbackNode = null),
							(e.callbackPriority = 0),
							Xu(Le, function () {
								return Uu(), null;
							}))
						: ((e.callbackNode = null), (e.callbackPriority = 0)),
					(r = (t.flags & 13878) != 0),
					t.subtreeFlags & 13878 || r)
				) {
					(r = T.T), (T.T = null), (i = E.p), (E.p = 2), (o = G), (G |= 4);
					try {
						ol(e, t, n);
					} finally {
						(G = o), (E.p = i), (T.T = r);
					}
				}
				(iu = 1), Ru(), zu(), Bu();
			}
		}
		function Ru() {
			if (iu === 1) {
				iu = 0;
				var e = au,
					t = ou,
					n = (t.flags & 13878) != 0;
				if (t.subtreeFlags & 13878 || n) {
					(n = T.T), (T.T = null);
					var r = E.p;
					E.p = 2;
					var i = G;
					G |= 4;
					try {
						vl(t, e);
						var a = zd,
							o = Nr(e.containerInfo),
							s = a.focusedElem,
							c = a.selectionRange;
						if (o !== s && s && s.ownerDocument && Mr(s.ownerDocument.documentElement, s)) {
							if (c !== null && Pr(s)) {
								var l = c.start,
									u = c.end;
								if ((u === void 0 && (u = l), `selectionStart` in s))
									(s.selectionStart = l), (s.selectionEnd = Math.min(u, s.value.length));
								else {
									var d = s.ownerDocument || document,
										f = (d && d.defaultView) || window;
									if (f.getSelection) {
										var p = f.getSelection(),
											m = s.textContent.length,
											h = Math.min(c.start, m),
											g = c.end === void 0 ? h : Math.min(c.end, m);
										!p.extend && h > g && ((o = g), (g = h), (h = o));
										var _ = jr(s, h),
											v = jr(s, g);
										if (
											_ &&
											v &&
											(p.rangeCount !== 1 ||
												p.anchorNode !== _.node ||
												p.anchorOffset !== _.offset ||
												p.focusNode !== v.node ||
												p.focusOffset !== v.offset)
										) {
											var y = d.createRange();
											y.setStart(_.node, _.offset),
												p.removeAllRanges(),
												h > g
													? (p.addRange(y), p.extend(v.node, v.offset))
													: (y.setEnd(v.node, v.offset), p.addRange(y));
										}
									}
								}
							}
							for (d = [], p = s; (p = p.parentNode); )
								p.nodeType === 1 && d.push({ element: p, left: p.scrollLeft, top: p.scrollTop });
							for (typeof s.focus == `function` && s.focus(), s = 0; s < d.length; s++) {
								var b = d[s];
								(b.element.scrollLeft = b.left), (b.element.scrollTop = b.top);
							}
						}
						(sp = !!Rd), (zd = Rd = null);
					} finally {
						(G = i), (E.p = r), (T.T = n);
					}
				}
				(e.current = t), (iu = 2);
			}
		}
		function zu() {
			if (iu === 2) {
				iu = 0;
				var e = au,
					t = ou,
					n = (t.flags & 8772) != 0;
				if (t.subtreeFlags & 8772 || n) {
					(n = T.T), (T.T = null);
					var r = E.p;
					E.p = 2;
					var i = G;
					G |= 4;
					try {
						sl(e, t.alternate, t);
					} finally {
						(G = i), (E.p = r), (T.T = n);
					}
				}
				iu = 3;
			}
		}
		function Bu() {
			if (iu === 4 || iu === 3) {
				(iu = 0), Me();
				var e = au,
					t = ou,
					n = su,
					r = uu;
				t.subtreeFlags & 10256 || t.flags & 10256
					? (iu = 5)
					: ((iu = 0), (ou = au = null), Vu(e, e.pendingLanes));
				var i = e.pendingLanes;
				if (
					(i === 0 && (ru = null),
					ut(n),
					(t = t.stateNode),
					Ue && typeof Ue.onCommitFiberRoot == `function`)
				)
					try {
						Ue.onCommitFiberRoot(He, t, void 0, (t.current.flags & 128) == 128);
					} catch {}
				if (r !== null) {
					(t = T.T), (i = E.p), (E.p = 2), (T.T = null);
					try {
						for (var a = e.onRecoverableError, o = 0; o < r.length; o++) {
							var s = r[o];
							a(s.value, { componentStack: s.stack });
						}
					} finally {
						(T.T = t), (E.p = i);
					}
				}
				su & 3 && Hu(),
					rd(e),
					(i = e.pendingLanes),
					n & 261930 && i & 42 ? (e === fu ? du++ : ((du = 0), (fu = e))) : (du = 0),
					id(0, !1);
			}
		}
		function Vu(e, t) {
			(e.pooledCacheLanes &= t) === 0 &&
				((t = e.pooledCache), t != null && ((e.pooledCache = null), pa(t)));
		}
		function Hu() {
			return Ru(), zu(), Bu(), Uu();
		}
		function Uu() {
			if (iu !== 5) return !1;
			var e = au,
				t = cu;
			cu = 0;
			var n = ut(su),
				r = T.T,
				i = E.p;
			try {
				(E.p = 32 > n ? 32 : n), (T.T = null), (n = lu), (lu = null);
				var a = au,
					o = su;
				if (((iu = 0), (ou = au = null), (su = 0), G & 6)) throw Error(s(331));
				var c = G;
				if (
					((G |= 4),
					Fl(a.current),
					Dl(a, a.current, o, n),
					(G = c),
					id(0, !1),
					Ue && typeof Ue.onPostCommitFiberRoot == `function`)
				)
					try {
						Ue.onPostCommitFiberRoot(He, a);
					} catch {}
				return !0;
			} finally {
				(E.p = i), (T.T = r), Vu(e, t);
			}
		}
		function Wu(e, t, n) {
			(t = wi(n, t)),
				(t = ec(e.stateNode, t, 2)),
				(e = Ja(e, t, 2)),
				e !== null && (it(e, 2), rd(e));
		}
		function Z(e, t, n) {
			if (e.tag === 3) Wu(e, e, n);
			else
				for (; t !== null; ) {
					if (t.tag === 3) {
						Wu(t, e, n);
						break;
					} else if (t.tag === 1) {
						var r = t.stateNode;
						if (
							typeof t.type.getDerivedStateFromError == `function` ||
							(typeof r.componentDidCatch == `function` && (ru === null || !ru.has(r)))
						) {
							(e = wi(n, e)),
								(n = tc(2)),
								(r = Ja(t, n, 2)),
								r !== null && (nc(n, r, t, e), it(r, 2), rd(r));
							break;
						}
					}
					t = t.return;
				}
		}
		function Gu(e, t, n) {
			var r = e.pingCache;
			if (r === null) {
				r = e.pingCache = new zl();
				var i = new Set();
				r.set(t, i);
			} else (i = r.get(t)), i === void 0 && ((i = new Set()), r.set(t, i));
			i.has(n) || ((Ul = !0), i.add(n), (e = Ku.bind(null, e, t, n)), t.then(e, e));
		}
		function Ku(e, t, n) {
			var r = e.pingCache;
			r !== null && r.delete(t),
				(e.pingedLanes |= e.suspendedLanes & n),
				(e.warmLanes &= ~n),
				K === e &&
					(J & n) === n &&
					(X === 4 || (X === 3 && (J & 62914560) === J && 300 > Ne() - $l)
						? !(G & 2) && Su(e, 0)
						: (ql |= n),
					Yl === J && (Yl = 0)),
				rd(e);
		}
		function qu(e, t) {
			t === 0 && (t = nt()), (e = li(e, t)), e !== null && (it(e, t), rd(e));
		}
		function Ju(e) {
			var t = e.memoizedState,
				n = 0;
			t !== null && (n = t.retryLane), qu(e, n);
		}
		function Yu(e, t) {
			var n = 0;
			switch (e.tag) {
				case 31:
				case 13:
					var r = e.stateNode,
						i = e.memoizedState;
					i !== null && (n = i.retryLane);
					break;
				case 19:
					r = e.stateNode;
					break;
				case 22:
					r = e.stateNode._retryCache;
					break;
				default:
					throw Error(s(314));
			}
			r !== null && r.delete(t), qu(e, n);
		}
		function Xu(e, t) {
			return ke(e, t);
		}
		var Zu = null,
			Qu = null,
			$u = !1,
			ed = !1,
			td = !1,
			nd = 0;
		function rd(e) {
			e !== Qu && e.next === null && (Qu === null ? (Zu = Qu = e) : (Qu = Qu.next = e)),
				(ed = !0),
				$u || (($u = !0), ud());
		}
		function id(e, t) {
			if (!td && ed) {
				td = !0;
				do
					for (var n = !1, r = Zu; r !== null; ) {
						if (!t)
							if (e !== 0) {
								var i = r.pendingLanes;
								if (i === 0) var a = 0;
								else {
									var o = r.suspendedLanes,
										s = r.pingedLanes;
									(a = (1 << (31 - Ge(42 | e) + 1)) - 1),
										(a &= i & ~(o & ~s)),
										(a = a & 201326741 ? (a & 201326741) | 1 : a ? a | 2 : 0);
								}
								a !== 0 && ((n = !0), ld(r, a));
							} else
								(a = J),
									(a = $e(
										r,
										r === K ? a : 0,
										r.cancelPendingCommit !== null || r.timeoutHandle !== -1,
									)),
									!(a & 3) || et(r, a) || ((n = !0), ld(r, a));
						r = r.next;
					}
				while (n);
				td = !1;
			}
		}
		function ad() {
			od();
		}
		function od() {
			ed = $u = !1;
			var e = 0;
			nd !== 0 && Gd() && (e = nd);
			for (var t = Ne(), n = null, r = Zu; r !== null; ) {
				var i = r.next,
					a = sd(r, t);
				a === 0
					? ((r.next = null), n === null ? (Zu = i) : (n.next = i), i === null && (Qu = n))
					: ((n = r), (e !== 0 || a & 3) && (ed = !0)),
					(r = i);
			}
			(iu !== 0 && iu !== 5) || id(e, !1), nd !== 0 && (nd = 0);
		}
		function sd(e, t) {
			for (
				var n = e.suspendedLanes,
					r = e.pingedLanes,
					i = e.expirationTimes,
					a = e.pendingLanes & -62914561;
				0 < a;
			) {
				var o = 31 - Ge(a),
					s = 1 << o,
					c = i[o];
				c === -1
					? ((s & n) === 0 || (s & r) !== 0) && (i[o] = tt(s, t))
					: c <= t && (e.expiredLanes |= s),
					(a &= ~s);
			}
			if (
				((t = K),
				(n = J),
				(n = $e(e, e === t ? n : 0, e.cancelPendingCommit !== null || e.timeoutHandle !== -1)),
				(r = e.callbackNode),
				n === 0 || (e === t && (Y === 2 || Y === 9)) || e.cancelPendingCommit !== null)
			)
				return r !== null && r !== null && Ae(r), (e.callbackNode = null), (e.callbackPriority = 0);
			if (!(n & 3) || et(e, n)) {
				if (((t = n & -n), t === e.callbackPriority)) return t;
				switch ((r !== null && Ae(r), ut(n))) {
					case 2:
					case 8:
						n = Ie;
						break;
					case 32:
						n = Le;
						break;
					case 268435456:
						n = ze;
						break;
					default:
						n = Le;
				}
				return (
					(r = cd.bind(null, e)), (n = ke(n, r)), (e.callbackPriority = t), (e.callbackNode = n), t
				);
			}
			return (
				r !== null && r !== null && Ae(r), (e.callbackPriority = 2), (e.callbackNode = null), 2
			);
		}
		function cd(e, t) {
			if (iu !== 0 && iu !== 5) return (e.callbackNode = null), (e.callbackPriority = 0), null;
			var n = e.callbackNode;
			if (Hu() && e.callbackNode !== n) return null;
			var r = J;
			return (
				(r = $e(e, e === K ? r : 0, e.cancelPendingCommit !== null || e.timeoutHandle !== -1)),
				r === 0
					? null
					: (gu(e, r, t),
						sd(e, Ne()),
						e.callbackNode != null && e.callbackNode === n ? cd.bind(null, e) : null)
			);
		}
		function ld(e, t) {
			if (Hu()) return null;
			gu(e, t, !0);
		}
		function ud() {
			Yd(function () {
				G & 6 ? ke(Fe, ad) : od();
			});
		}
		function dd() {
			if (nd === 0) {
				var e = ga;
				e === 0 && ((e = Ye), (Ye <<= 1), !(Ye & 261888) && (Ye = 256)), (nd = e);
			}
			return nd;
		}
		function fd(e) {
			return e == null || typeof e == `symbol` || typeof e == `boolean`
				? null
				: typeof e == `function`
					? e
					: on(`` + e);
		}
		function pd(e, t) {
			var n = t.ownerDocument.createElement(`input`);
			return (
				(n.name = t.name),
				(n.value = t.value),
				e.id && n.setAttribute(`form`, e.id),
				t.parentNode.insertBefore(n, t),
				(e = new FormData(e)),
				n.parentNode.removeChild(n),
				e
			);
		}
		function md(e, t, n, r, i) {
			if (t === `submit` && n && n.stateNode === i) {
				var a = fd((i[ht] || null).action),
					o = r.submitter;
				o &&
					((t = (t = o[ht] || null) ? fd(t.formAction) : o.getAttribute(`formAction`)),
					t !== null && ((a = t), (o = null)));
				var s = new On(`action`, `action`, null, r, i);
				e.push({
					event: s,
					listeners: [
						{
							instance: null,
							listener: function () {
								if (r.defaultPrevented) {
									if (nd !== 0) {
										var e = o ? pd(i, o) : new FormData(i);
										Es(n, { pending: !0, data: e, method: i.method, action: a }, null, e);
									}
								} else
									typeof a == `function` &&
										(s.preventDefault(),
										(e = o ? pd(i, o) : new FormData(i)),
										Es(n, { pending: !0, data: e, method: i.method, action: a }, a, e));
							},
							currentTarget: i,
						},
					],
				});
			}
		}
		for (var hd = 0; hd < ei.length; hd++) {
			var gd = ei[hd];
			ti(gd.toLowerCase(), `on` + (gd[0].toUpperCase() + gd.slice(1)));
		}
		ti(Kr, `onAnimationEnd`),
			ti(qr, `onAnimationIteration`),
			ti(Jr, `onAnimationStart`),
			ti(`dblclick`, `onDoubleClick`),
			ti(`focusin`, `onFocus`),
			ti(`focusout`, `onBlur`),
			ti(Yr, `onTransitionRun`),
			ti(Xr, `onTransitionStart`),
			ti(Zr, `onTransitionCancel`),
			ti(Qr, `onTransitionEnd`),
			jt(`onMouseEnter`, [`mouseout`, `mouseover`]),
			jt(`onMouseLeave`, [`mouseout`, `mouseover`]),
			jt(`onPointerEnter`, [`pointerout`, `pointerover`]),
			jt(`onPointerLeave`, [`pointerout`, `pointerover`]),
			At(
				`onChange`,
				`change click focusin focusout input keydown keyup selectionchange`.split(` `),
			),
			At(
				`onSelect`,
				`focusout contextmenu dragend focusin keydown keyup mousedown mouseup selectionchange`.split(
					` `,
				),
			),
			At(`onBeforeInput`, [`compositionend`, `keypress`, `textInput`, `paste`]),
			At(`onCompositionEnd`, `compositionend focusout keydown keypress keyup mousedown`.split(` `)),
			At(
				`onCompositionStart`,
				`compositionstart focusout keydown keypress keyup mousedown`.split(` `),
			),
			At(
				`onCompositionUpdate`,
				`compositionupdate focusout keydown keypress keyup mousedown`.split(` `),
			);
		var _d =
				`abort canplay canplaythrough durationchange emptied encrypted ended error loadeddata loadedmetadata loadstart pause play playing progress ratechange resize seeked seeking stalled suspend timeupdate volumechange waiting`.split(
					` `,
				),
			vd = new Set(
				`beforetoggle cancel close invalid load scroll scrollend toggle`.split(` `).concat(_d),
			);
		function yd(e, t) {
			t = (t & 4) != 0;
			for (var n = 0; n < e.length; n++) {
				var r = e[n],
					i = r.event;
				r = r.listeners;
				a: {
					var a = void 0;
					if (t)
						for (var o = r.length - 1; 0 <= o; o--) {
							var s = r[o],
								c = s.instance,
								l = s.currentTarget;
							if (((s = s.listener), c !== a && i.isPropagationStopped())) break a;
							(a = s), (i.currentTarget = l);
							try {
								a(i);
							} catch (e) {
								ni(e);
							}
							(i.currentTarget = null), (a = c);
						}
					else
						for (o = 0; o < r.length; o++) {
							if (
								((s = r[o]),
								(c = s.instance),
								(l = s.currentTarget),
								(s = s.listener),
								c !== a && i.isPropagationStopped())
							)
								break a;
							(a = s), (i.currentTarget = l);
							try {
								a(i);
							} catch (e) {
								ni(e);
							}
							(i.currentTarget = null), (a = c);
						}
				}
			}
		}
		function Q(e, t) {
			var n = t[_t];
			n === void 0 && (n = t[_t] = new Set());
			var r = e + `__bubble`;
			n.has(r) || (Cd(t, e, 2, !1), n.add(r));
		}
		function bd(e, t, n) {
			var r = 0;
			t && (r |= 4), Cd(n, e, r, t);
		}
		var xd = `_reactListening` + Math.random().toString(36).slice(2);
		function Sd(e) {
			if (!e[xd]) {
				(e[xd] = !0),
					Ot.forEach(function (t) {
						t !== `selectionchange` && (vd.has(t) || bd(t, !1, e), bd(t, !0, e));
					});
				var t = e.nodeType === 9 ? e : e.ownerDocument;
				t === null || t[xd] || ((t[xd] = !0), bd(`selectionchange`, !1, t));
			}
		}
		function Cd(e, t, n, r) {
			switch (mp(t)) {
				case 2:
					var i = cp;
					break;
				case 8:
					i = lp;
					break;
				default:
					i = up;
			}
			(n = i.bind(null, t, n, e)),
				(i = void 0),
				!_n || (t !== `touchstart` && t !== `touchmove` && t !== `wheel`) || (i = !0),
				r
					? i === void 0
						? e.addEventListener(t, n, !0)
						: e.addEventListener(t, n, { capture: !0, passive: i })
					: i === void 0
						? e.addEventListener(t, n, !1)
						: e.addEventListener(t, n, { passive: i });
		}
		function wd(e, t, n, r, i) {
			var a = r;
			if (!(t & 1) && !(t & 2) && r !== null)
				a: for (;;) {
					if (r === null) return;
					var o = r.tag;
					if (o === 3 || o === 4) {
						var s = r.stateNode.containerInfo;
						if (s === i) break;
						if (o === 4)
							for (o = r.return; o !== null; ) {
								var c = o.tag;
								if ((c === 3 || c === 4) && o.stateNode.containerInfo === i) return;
								o = o.return;
							}
						for (; s !== null; ) {
							if (((o = Ct(s)), o === null)) return;
							if (((c = o.tag), c === 5 || c === 6 || c === 26 || c === 27)) {
								r = a = o;
								continue a;
							}
							s = s.parentNode;
						}
					}
					r = r.return;
				}
			mn(function () {
				var r = a,
					i = ln(n),
					o = [];
				a: {
					var s = $r.get(e);
					if (s !== void 0) {
						var c = On,
							u = e;
						switch (e) {
							case `keypress`:
								if (Cn(n) === 0) break a;
							case `keydown`:
							case `keyup`:
								c = Kn;
								break;
							case `focusin`:
								(u = `focus`), (c = Ln);
								break;
							case `focusout`:
								(u = `blur`), (c = Ln);
								break;
							case `beforeblur`:
							case `afterblur`:
								c = Ln;
								break;
							case `click`:
								if (n.button === 2) break a;
							case `auxclick`:
							case `dblclick`:
							case `mousedown`:
							case `mousemove`:
							case `mouseup`:
							case `mouseout`:
							case `mouseover`:
							case `contextmenu`:
								c = Fn;
								break;
							case `drag`:
							case `dragend`:
							case `dragenter`:
							case `dragexit`:
							case `dragleave`:
							case `dragover`:
							case `dragstart`:
							case `drop`:
								c = In;
								break;
							case `touchcancel`:
							case `touchend`:
							case `touchmove`:
							case `touchstart`:
								c = Jn;
								break;
							case Kr:
							case qr:
							case Jr:
								c = Rn;
								break;
							case Qr:
								c = M;
								break;
							case `scroll`:
							case `scrollend`:
								c = An;
								break;
							case `wheel`:
								c = Yn;
								break;
							case `copy`:
							case `cut`:
							case `paste`:
								c = zn;
								break;
							case `gotpointercapture`:
							case `lostpointercapture`:
							case `pointercancel`:
							case `pointerdown`:
							case `pointermove`:
							case `pointerout`:
							case `pointerover`:
							case `pointerup`:
								c = qn;
								break;
							case `toggle`:
							case `beforetoggle`:
								c = Xn;
						}
						var d = (t & 4) != 0,
							f = !d && (e === `scroll` || e === `scrollend`),
							p = d ? (s === null ? null : s + `Capture`) : s;
						d = [];
						for (var m = r, h; m !== null; ) {
							var g = m;
							if (
								((h = g.stateNode),
								(g = g.tag),
								(g !== 5 && g !== 26 && g !== 27) ||
									h === null ||
									p === null ||
									((g = hn(m, p)), g != null && d.push(Td(m, g, h))),
								f)
							)
								break;
							m = m.return;
						}
						0 < d.length && ((s = new c(s, u, null, n, i)), o.push({ event: s, listeners: d }));
					}
				}
				if (!(t & 7)) {
					a: {
						if (
							((s = e === `mouseover` || e === `pointerover`),
							(c = e === `mouseout` || e === `pointerout`),
							s && n !== cn && (u = n.relatedTarget || n.fromElement) && (Ct(u) || u[gt]))
						)
							break a;
						if (
							(c || s) &&
							((s =
								i.window === i
									? i
									: (s = i.ownerDocument)
										? s.defaultView || s.parentWindow
										: window),
							c
								? ((u = n.relatedTarget || n.toElement),
									(c = r),
									(u = u ? Ct(u) : null),
									u !== null &&
										((f = l(u)), (d = u.tag), u !== f || (d !== 5 && d !== 27 && d !== 6)) &&
										(u = null))
								: ((c = null), (u = r)),
							c !== u)
						) {
							if (
								((d = Fn),
								(g = `onMouseLeave`),
								(p = `onMouseEnter`),
								(m = `mouse`),
								(e === `pointerout` || e === `pointerover`) &&
									((d = qn), (g = `onPointerLeave`), (p = `onPointerEnter`), (m = `pointer`)),
								(f = c == null ? s : Tt(c)),
								(h = u == null ? s : Tt(u)),
								(s = new d(g, m + `leave`, c, n, i)),
								(s.target = f),
								(s.relatedTarget = h),
								(g = null),
								Ct(i) === r &&
									((d = new d(p, m + `enter`, u, n, i)),
									(d.target = h),
									(d.relatedTarget = f),
									(g = d)),
								(f = g),
								c && u)
							)
								b: {
									for (d = Dd, p = c, m = u, h = 0, g = p; g; g = d(g)) h++;
									g = 0;
									for (var _ = m; _; _ = d(_)) g++;
									for (; 0 < h - g; ) (p = d(p)), h--;
									for (; 0 < g - h; ) (m = d(m)), g--;
									for (; h--; ) {
										if (p === m || (m !== null && p === m.alternate)) {
											d = p;
											break b;
										}
										(p = d(p)), (m = d(m));
									}
									d = null;
								}
							else d = null;
							c !== null && Od(o, s, c, d, !1), u !== null && f !== null && Od(o, f, u, d, !0);
						}
					}
					a: {
						if (
							((s = r ? Tt(r) : window),
							(c = s.nodeName && s.nodeName.toLowerCase()),
							c === `select` || (c === `input` && s.type === `file`))
						)
							var v = gr;
						else if (ur(s))
							if (_r) v = Er;
							else {
								v = wr;
								var y = Cr;
							}
						else
							(c = s.nodeName),
								!c || c.toLowerCase() !== `input` || (s.type !== `checkbox` && s.type !== `radio`)
									? r && nn(r.elementType) && (v = gr)
									: (v = Tr);
						if ((v &&= v(e, r))) {
							dr(o, v, n, i);
							break a;
						}
						y && y(e, s, r),
							e === `focusout` &&
								r &&
								s.type === `number` &&
								r.memoizedProps.value != null &&
								Yt(s, `number`, s.value);
					}
					switch (((y = r ? Tt(r) : window), e)) {
						case `focusin`:
							(ur(y) || y.contentEditable === `true`) && ((Ir = y), (Lr = r), (Rr = null));
							break;
						case `focusout`:
							Rr = Lr = Ir = null;
							break;
						case `mousedown`:
							zr = !0;
							break;
						case `contextmenu`:
						case `mouseup`:
						case `dragend`:
							(zr = !1), Br(o, n, i);
							break;
						case `selectionchange`:
							if (Fr) break;
						case `keydown`:
						case `keyup`:
							Br(o, n, i);
					}
					var b;
					if (Qn)
						b: {
							switch (e) {
								case `compositionstart`:
									var x = `onCompositionStart`;
									break b;
								case `compositionend`:
									x = `onCompositionEnd`;
									break b;
								case `compositionupdate`:
									x = `onCompositionUpdate`;
									break b;
							}
							x = void 0;
						}
					else
						or
							? ir(e, n) && (x = `onCompositionEnd`)
							: e === `keydown` && n.keyCode === 229 && (x = `onCompositionStart`);
					x &&
						(tr &&
							n.locale !== `ko` &&
							(or || x !== `onCompositionStart`
								? x === `onCompositionEnd` && or && (b = Sn())
								: ((yn = i), (bn = `value` in yn ? yn.value : yn.textContent), (or = !0))),
						(y = Ed(r, x)),
						0 < y.length &&
							((x = new Bn(x, e, null, n, i)),
							o.push({ event: x, listeners: y }),
							b ? (x.data = b) : ((b = ar(n)), b !== null && (x.data = b)))),
						(b = er ? sr(e, n) : cr(e, n)) &&
							((x = Ed(r, `onBeforeInput`)),
							0 < x.length &&
								((y = new Bn(`onBeforeInput`, `beforeinput`, null, n, i)),
								o.push({ event: y, listeners: x }),
								(y.data = b))),
						md(o, e, r, n, i);
				}
				yd(o, t);
			});
		}
		function Td(e, t, n) {
			return { instance: e, listener: t, currentTarget: n };
		}
		function Ed(e, t) {
			for (var n = t + `Capture`, r = []; e !== null; ) {
				var i = e,
					a = i.stateNode;
				if (
					((i = i.tag),
					(i !== 5 && i !== 26 && i !== 27) ||
						a === null ||
						((i = hn(e, n)),
						i != null && r.unshift(Td(e, i, a)),
						(i = hn(e, t)),
						i != null && r.push(Td(e, i, a))),
					e.tag === 3)
				)
					return r;
				e = e.return;
			}
			return [];
		}
		function Dd(e) {
			if (e === null) return null;
			do e = e.return;
			while (e && e.tag !== 5 && e.tag !== 27);
			return e || null;
		}
		function Od(e, t, n, r, i) {
			for (var a = t._reactName, o = []; n !== null && n !== r; ) {
				var s = n,
					c = s.alternate,
					l = s.stateNode;
				if (((s = s.tag), c !== null && c === r)) break;
				(s !== 5 && s !== 26 && s !== 27) ||
					l === null ||
					((c = l),
					i
						? ((l = hn(n, a)), l != null && o.unshift(Td(n, l, c)))
						: i || ((l = hn(n, a)), l != null && o.push(Td(n, l, c)))),
					(n = n.return);
			}
			o.length !== 0 && e.push({ event: t, listeners: o });
		}
		var kd = /\r\n?/g,
			Ad = /\u0000|\uFFFD/g;
		function jd(e) {
			return (typeof e == `string` ? e : `` + e)
				.replace(
					kd,
					`
`,
				)
				.replace(Ad, ``);
		}
		function Md(e, t) {
			return (t = jd(t)), jd(e) === t;
		}
		function $(e, t, n, r, i, a) {
			switch (n) {
				case `children`:
					typeof r == `string`
						? t === `body` || (t === `textarea` && r === ``) || j(e, r)
						: (typeof r == `number` || typeof r == `bigint`) && t !== `body` && j(e, `` + r);
					break;
				case `className`:
					Lt(e, `class`, r);
					break;
				case `tabIndex`:
					Lt(e, `tabindex`, r);
					break;
				case `dir`:
				case `role`:
				case `viewBox`:
				case `width`:
				case `height`:
					Lt(e, n, r);
					break;
				case `style`:
					tn(e, r, a);
					break;
				case `data`:
					if (t !== `object`) {
						Lt(e, `data`, r);
						break;
					}
				case `src`:
				case `href`:
					if (r === `` && (t !== `a` || n !== `href`)) {
						e.removeAttribute(n);
						break;
					}
					if (
						r == null ||
						typeof r == `function` ||
						typeof r == `symbol` ||
						typeof r == `boolean`
					) {
						e.removeAttribute(n);
						break;
					}
					(r = on(`` + r)), e.setAttribute(n, r);
					break;
				case `action`:
				case `formAction`:
					if (typeof r == `function`) {
						e.setAttribute(
							n,
							`javascript:throw new Error('A React form was unexpectedly submitted. If you called form.submit() manually, consider using form.requestSubmit() instead. If you\\'re trying to use event.stopPropagation() in a submit event handler, consider also calling event.preventDefault().')`,
						);
						break;
					} else
						typeof a == `function` &&
							(n === `formAction`
								? (t !== `input` && $(e, t, `name`, i.name, i, null),
									$(e, t, `formEncType`, i.formEncType, i, null),
									$(e, t, `formMethod`, i.formMethod, i, null),
									$(e, t, `formTarget`, i.formTarget, i, null))
								: ($(e, t, `encType`, i.encType, i, null),
									$(e, t, `method`, i.method, i, null),
									$(e, t, `target`, i.target, i, null)));
					if (r == null || typeof r == `symbol` || typeof r == `boolean`) {
						e.removeAttribute(n);
						break;
					}
					(r = on(`` + r)), e.setAttribute(n, r);
					break;
				case `onClick`:
					r != null && (e.onclick = sn);
					break;
				case `onScroll`:
					r != null && Q(`scroll`, e);
					break;
				case `onScrollEnd`:
					r != null && Q(`scrollend`, e);
					break;
				case `dangerouslySetInnerHTML`:
					if (r != null) {
						if (typeof r != `object` || !(`__html` in r)) throw Error(s(61));
						if (((n = r.__html), n != null)) {
							if (i.children != null) throw Error(s(60));
							e.innerHTML = n;
						}
					}
					break;
				case `multiple`:
					e.multiple = r && typeof r != `function` && typeof r != `symbol`;
					break;
				case `muted`:
					e.muted = r && typeof r != `function` && typeof r != `symbol`;
					break;
				case `suppressContentEditableWarning`:
				case `suppressHydrationWarning`:
				case `defaultValue`:
				case `defaultChecked`:
				case `innerHTML`:
				case `ref`:
					break;
				case `autoFocus`:
					break;
				case `xlinkHref`:
					if (
						r == null ||
						typeof r == `function` ||
						typeof r == `boolean` ||
						typeof r == `symbol`
					) {
						e.removeAttribute(`xlink:href`);
						break;
					}
					(n = on(`` + r)), e.setAttributeNS(`http://www.w3.org/1999/xlink`, `xlink:href`, n);
					break;
				case `contentEditable`:
				case `spellCheck`:
				case `draggable`:
				case `value`:
				case `autoReverse`:
				case `externalResourcesRequired`:
				case `focusable`:
				case `preserveAlpha`:
					r != null && typeof r != `function` && typeof r != `symbol`
						? e.setAttribute(n, `` + r)
						: e.removeAttribute(n);
					break;
				case `inert`:
				case `allowFullScreen`:
				case `async`:
				case `autoPlay`:
				case `controls`:
				case `default`:
				case `defer`:
				case `disabled`:
				case `disablePictureInPicture`:
				case `disableRemotePlayback`:
				case `formNoValidate`:
				case `hidden`:
				case `loop`:
				case `noModule`:
				case `noValidate`:
				case `open`:
				case `playsInline`:
				case `readOnly`:
				case `required`:
				case `reversed`:
				case `scoped`:
				case `seamless`:
				case `itemScope`:
					r && typeof r != `function` && typeof r != `symbol`
						? e.setAttribute(n, ``)
						: e.removeAttribute(n);
					break;
				case `capture`:
				case `download`:
					!0 === r
						? e.setAttribute(n, ``)
						: !1 !== r && r != null && typeof r != `function` && typeof r != `symbol`
							? e.setAttribute(n, r)
							: e.removeAttribute(n);
					break;
				case `cols`:
				case `rows`:
				case `size`:
				case `span`:
					r != null && typeof r != `function` && typeof r != `symbol` && !isNaN(r) && 1 <= r
						? e.setAttribute(n, r)
						: e.removeAttribute(n);
					break;
				case `rowSpan`:
				case `start`:
					r == null || typeof r == `function` || typeof r == `symbol` || isNaN(r)
						? e.removeAttribute(n)
						: e.setAttribute(n, r);
					break;
				case `popover`:
					Q(`beforetoggle`, e), Q(`toggle`, e), It(e, `popover`, r);
					break;
				case `xlinkActuate`:
					Rt(e, `http://www.w3.org/1999/xlink`, `xlink:actuate`, r);
					break;
				case `xlinkArcrole`:
					Rt(e, `http://www.w3.org/1999/xlink`, `xlink:arcrole`, r);
					break;
				case `xlinkRole`:
					Rt(e, `http://www.w3.org/1999/xlink`, `xlink:role`, r);
					break;
				case `xlinkShow`:
					Rt(e, `http://www.w3.org/1999/xlink`, `xlink:show`, r);
					break;
				case `xlinkTitle`:
					Rt(e, `http://www.w3.org/1999/xlink`, `xlink:title`, r);
					break;
				case `xlinkType`:
					Rt(e, `http://www.w3.org/1999/xlink`, `xlink:type`, r);
					break;
				case `xmlBase`:
					Rt(e, `http://www.w3.org/XML/1998/namespace`, `xml:base`, r);
					break;
				case `xmlLang`:
					Rt(e, `http://www.w3.org/XML/1998/namespace`, `xml:lang`, r);
					break;
				case `xmlSpace`:
					Rt(e, `http://www.w3.org/XML/1998/namespace`, `xml:space`, r);
					break;
				case `is`:
					It(e, `is`, r);
					break;
				case `innerText`:
				case `textContent`:
					break;
				default:
					(!(2 < n.length) || (n[0] !== `o` && n[0] !== `O`) || (n[1] !== `n` && n[1] !== `N`)) &&
						((n = rn.get(n) || n), It(e, n, r));
			}
		}
		function Nd(e, t, n, r, i, a) {
			switch (n) {
				case `style`:
					tn(e, r, a);
					break;
				case `dangerouslySetInnerHTML`:
					if (r != null) {
						if (typeof r != `object` || !(`__html` in r)) throw Error(s(61));
						if (((n = r.__html), n != null)) {
							if (i.children != null) throw Error(s(60));
							e.innerHTML = n;
						}
					}
					break;
				case `children`:
					typeof r == `string`
						? j(e, r)
						: (typeof r == `number` || typeof r == `bigint`) && j(e, `` + r);
					break;
				case `onScroll`:
					r != null && Q(`scroll`, e);
					break;
				case `onScrollEnd`:
					r != null && Q(`scrollend`, e);
					break;
				case `onClick`:
					r != null && (e.onclick = sn);
					break;
				case `suppressContentEditableWarning`:
				case `suppressHydrationWarning`:
				case `innerHTML`:
				case `ref`:
					break;
				case `innerText`:
				case `textContent`:
					break;
				default:
					if (!kt.hasOwnProperty(n))
						a: {
							if (
								n[0] === `o` &&
								n[1] === `n` &&
								((i = n.endsWith(`Capture`)),
								(t = n.slice(2, i ? n.length - 7 : void 0)),
								(a = e[ht] || null),
								(a = a == null ? null : a[n]),
								typeof a == `function` && e.removeEventListener(t, a, i),
								typeof r == `function`)
							) {
								typeof a != `function` &&
									a !== null &&
									(n in e ? (e[n] = null) : e.hasAttribute(n) && e.removeAttribute(n)),
									e.addEventListener(t, r, i);
								break a;
							}
							n in e ? (e[n] = r) : !0 === r ? e.setAttribute(n, ``) : It(e, n, r);
						}
			}
		}
		function Pd(e, t, n) {
			switch (t) {
				case `div`:
				case `span`:
				case `svg`:
				case `path`:
				case `a`:
				case `g`:
				case `p`:
				case `li`:
					break;
				case `img`:
					Q(`error`, e), Q(`load`, e);
					var r = !1,
						i = !1,
						a;
					for (a in n)
						if (n.hasOwnProperty(a)) {
							var o = n[a];
							if (o != null)
								switch (a) {
									case `src`:
										r = !0;
										break;
									case `srcSet`:
										i = !0;
										break;
									case `children`:
									case `dangerouslySetInnerHTML`:
										throw Error(s(137, t));
									default:
										$(e, t, a, o, n, null);
								}
						}
					i && $(e, t, `srcSet`, n.srcSet, n, null), r && $(e, t, `src`, n.src, n, null);
					return;
				case `input`:
					Q(`invalid`, e);
					var c = (a = o = i = null),
						l = null,
						u = null;
					for (r in n)
						if (n.hasOwnProperty(r)) {
							var d = n[r];
							if (d != null)
								switch (r) {
									case `name`:
										i = d;
										break;
									case `type`:
										o = d;
										break;
									case `checked`:
										l = d;
										break;
									case `defaultChecked`:
										u = d;
										break;
									case `value`:
										a = d;
										break;
									case `defaultValue`:
										c = d;
										break;
									case `children`:
									case `dangerouslySetInnerHTML`:
										if (d != null) throw Error(s(137, t));
										break;
									default:
										$(e, t, r, d, n, null);
								}
						}
					Jt(e, a, c, l, u, o, i, !1);
					return;
				case `select`:
					for (i in (Q(`invalid`, e), (r = o = a = null), n))
						if (n.hasOwnProperty(i) && ((c = n[i]), c != null))
							switch (i) {
								case `value`:
									a = c;
									break;
								case `defaultValue`:
									o = c;
									break;
								case `multiple`:
									r = c;
								default:
									$(e, t, i, c, n, null);
							}
					(t = a),
						(n = o),
						(e.multiple = !!r),
						t == null ? n != null && Xt(e, !!r, n, !0) : Xt(e, !!r, t, !1);
					return;
				case `textarea`:
					for (o in (Q(`invalid`, e), (a = i = r = null), n))
						if (n.hasOwnProperty(o) && ((c = n[o]), c != null))
							switch (o) {
								case `value`:
									r = c;
									break;
								case `defaultValue`:
									i = c;
									break;
								case `children`:
									a = c;
									break;
								case `dangerouslySetInnerHTML`:
									if (c != null) throw Error(s(91));
									break;
								default:
									$(e, t, o, c, n, null);
							}
					Qt(e, r, i, a);
					return;
				case `option`:
					for (l in n)
						if (n.hasOwnProperty(l) && ((r = n[l]), r != null))
							switch (l) {
								case `selected`:
									e.selected = r && typeof r != `function` && typeof r != `symbol`;
									break;
								default:
									$(e, t, l, r, n, null);
							}
					return;
				case `dialog`:
					Q(`beforetoggle`, e), Q(`toggle`, e), Q(`cancel`, e), Q(`close`, e);
					break;
				case `iframe`:
				case `object`:
					Q(`load`, e);
					break;
				case `video`:
				case `audio`:
					for (r = 0; r < _d.length; r++) Q(_d[r], e);
					break;
				case `image`:
					Q(`error`, e), Q(`load`, e);
					break;
				case `details`:
					Q(`toggle`, e);
					break;
				case `embed`:
				case `source`:
				case `link`:
					Q(`error`, e), Q(`load`, e);
				case `area`:
				case `base`:
				case `br`:
				case `col`:
				case `hr`:
				case `keygen`:
				case `meta`:
				case `param`:
				case `track`:
				case `wbr`:
				case `menuitem`:
					for (u in n)
						if (n.hasOwnProperty(u) && ((r = n[u]), r != null))
							switch (u) {
								case `children`:
								case `dangerouslySetInnerHTML`:
									throw Error(s(137, t));
								default:
									$(e, t, u, r, n, null);
							}
					return;
				default:
					if (nn(t)) {
						for (d in n)
							n.hasOwnProperty(d) && ((r = n[d]), r !== void 0 && Nd(e, t, d, r, n, void 0));
						return;
					}
			}
			for (c in n) n.hasOwnProperty(c) && ((r = n[c]), r != null && $(e, t, c, r, n, null));
		}
		function Fd(e, t, n, r) {
			switch (t) {
				case `div`:
				case `span`:
				case `svg`:
				case `path`:
				case `a`:
				case `g`:
				case `p`:
				case `li`:
					break;
				case `input`:
					var i = null,
						a = null,
						o = null,
						c = null,
						l = null,
						u = null,
						d = null;
					for (m in n) {
						var f = n[m];
						if (n.hasOwnProperty(m) && f != null)
							switch (m) {
								case `checked`:
									break;
								case `value`:
									break;
								case `defaultValue`:
									l = f;
								default:
									r.hasOwnProperty(m) || $(e, t, m, null, r, f);
							}
					}
					for (var p in r) {
						var m = r[p];
						if (((f = n[p]), r.hasOwnProperty(p) && (m != null || f != null)))
							switch (p) {
								case `type`:
									a = m;
									break;
								case `name`:
									i = m;
									break;
								case `checked`:
									u = m;
									break;
								case `defaultChecked`:
									d = m;
									break;
								case `value`:
									o = m;
									break;
								case `defaultValue`:
									c = m;
									break;
								case `children`:
								case `dangerouslySetInnerHTML`:
									if (m != null) throw Error(s(137, t));
									break;
								default:
									m !== f && $(e, t, p, m, r, f);
							}
					}
					qt(e, o, c, l, u, d, a, i);
					return;
				case `select`:
					for (a in ((m = o = c = p = null), n))
						if (((l = n[a]), n.hasOwnProperty(a) && l != null))
							switch (a) {
								case `value`:
									break;
								case `multiple`:
									m = l;
								default:
									r.hasOwnProperty(a) || $(e, t, a, null, r, l);
							}
					for (i in r)
						if (((a = r[i]), (l = n[i]), r.hasOwnProperty(i) && (a != null || l != null)))
							switch (i) {
								case `value`:
									p = a;
									break;
								case `defaultValue`:
									c = a;
									break;
								case `multiple`:
									o = a;
								default:
									a !== l && $(e, t, i, a, r, l);
							}
					(t = c),
						(n = o),
						(r = m),
						p == null
							? !!r != !!n && (t == null ? Xt(e, !!n, n ? [] : ``, !1) : Xt(e, !!n, t, !0))
							: Xt(e, !!n, p, !1);
					return;
				case `textarea`:
					for (c in ((m = p = null), n))
						if (((i = n[c]), n.hasOwnProperty(c) && i != null && !r.hasOwnProperty(c)))
							switch (c) {
								case `value`:
									break;
								case `children`:
									break;
								default:
									$(e, t, c, null, r, i);
							}
					for (o in r)
						if (((i = r[o]), (a = n[o]), r.hasOwnProperty(o) && (i != null || a != null)))
							switch (o) {
								case `value`:
									p = i;
									break;
								case `defaultValue`:
									m = i;
									break;
								case `children`:
									break;
								case `dangerouslySetInnerHTML`:
									if (i != null) throw Error(s(91));
									break;
								default:
									i !== a && $(e, t, o, i, r, a);
							}
					Zt(e, p, m);
					return;
				case `option`:
					for (var h in n)
						if (((p = n[h]), n.hasOwnProperty(h) && p != null && !r.hasOwnProperty(h)))
							switch (h) {
								case `selected`:
									e.selected = !1;
									break;
								default:
									$(e, t, h, null, r, p);
							}
					for (l in r)
						if (
							((p = r[l]), (m = n[l]), r.hasOwnProperty(l) && p !== m && (p != null || m != null))
						)
							switch (l) {
								case `selected`:
									e.selected = p && typeof p != `function` && typeof p != `symbol`;
									break;
								default:
									$(e, t, l, p, r, m);
							}
					return;
				case `img`:
				case `link`:
				case `area`:
				case `base`:
				case `br`:
				case `col`:
				case `embed`:
				case `hr`:
				case `keygen`:
				case `meta`:
				case `param`:
				case `source`:
				case `track`:
				case `wbr`:
				case `menuitem`:
					for (var g in n)
						(p = n[g]),
							n.hasOwnProperty(g) && p != null && !r.hasOwnProperty(g) && $(e, t, g, null, r, p);
					for (u in r)
						if (
							((p = r[u]), (m = n[u]), r.hasOwnProperty(u) && p !== m && (p != null || m != null))
						)
							switch (u) {
								case `children`:
								case `dangerouslySetInnerHTML`:
									if (p != null) throw Error(s(137, t));
									break;
								default:
									$(e, t, u, p, r, m);
							}
					return;
				default:
					if (nn(t)) {
						for (var _ in n)
							(p = n[_]),
								n.hasOwnProperty(_) &&
									p !== void 0 &&
									!r.hasOwnProperty(_) &&
									Nd(e, t, _, void 0, r, p);
						for (d in r)
							(p = r[d]),
								(m = n[d]),
								!r.hasOwnProperty(d) ||
									p === m ||
									(p === void 0 && m === void 0) ||
									Nd(e, t, d, p, r, m);
						return;
					}
			}
			for (var v in n)
				(p = n[v]),
					n.hasOwnProperty(v) && p != null && !r.hasOwnProperty(v) && $(e, t, v, null, r, p);
			for (f in r)
				(p = r[f]),
					(m = n[f]),
					!r.hasOwnProperty(f) || p === m || (p == null && m == null) || $(e, t, f, p, r, m);
		}
		function Id(e) {
			switch (e) {
				case `css`:
				case `script`:
				case `font`:
				case `img`:
				case `image`:
				case `input`:
				case `link`:
					return !0;
				default:
					return !1;
			}
		}
		function Ld() {
			if (typeof performance.getEntriesByType == `function`) {
				for (
					var e = 0, t = 0, n = performance.getEntriesByType(`resource`), r = 0;
					r < n.length;
					r++
				) {
					var i = n[r],
						a = i.transferSize,
						o = i.initiatorType,
						s = i.duration;
					if (a && s && Id(o)) {
						for (o = 0, s = i.responseEnd, r += 1; r < n.length; r++) {
							var c = n[r],
								l = c.startTime;
							if (l > s) break;
							var u = c.transferSize,
								d = c.initiatorType;
							u && Id(d) && ((c = c.responseEnd), (o += u * (c < s ? 1 : (s - l) / (c - l))));
						}
						if ((--r, (t += (8 * (a + o)) / (i.duration / 1e3)), e++, 10 < e)) break;
					}
				}
				if (0 < e) return t / e / 1e6;
			}
			return navigator.connection && ((e = navigator.connection.downlink), typeof e == `number`)
				? e
				: 5;
		}
		var Rd = null,
			zd = null;
		function Bd(e) {
			return e.nodeType === 9 ? e : e.ownerDocument;
		}
		function Vd(e) {
			switch (e) {
				case `http://www.w3.org/2000/svg`:
					return 1;
				case `http://www.w3.org/1998/Math/MathML`:
					return 2;
				default:
					return 0;
			}
		}
		function Hd(e, t) {
			if (e === 0)
				switch (t) {
					case `svg`:
						return 1;
					case `math`:
						return 2;
					default:
						return 0;
				}
			return e === 1 && t === `foreignObject` ? 0 : e;
		}
		function Ud(e, t) {
			return (
				e === `textarea` ||
				e === `noscript` ||
				typeof t.children == `string` ||
				typeof t.children == `number` ||
				typeof t.children == `bigint` ||
				(typeof t.dangerouslySetInnerHTML == `object` &&
					t.dangerouslySetInnerHTML !== null &&
					t.dangerouslySetInnerHTML.__html != null)
			);
		}
		var Wd = null;
		function Gd() {
			var e = window.event;
			return e && e.type === `popstate` ? (e === Wd ? !1 : ((Wd = e), !0)) : ((Wd = null), !1);
		}
		var Kd = typeof setTimeout == `function` ? setTimeout : void 0,
			qd = typeof clearTimeout == `function` ? clearTimeout : void 0,
			Jd = typeof Promise == `function` ? Promise : void 0,
			Yd =
				typeof queueMicrotask == `function`
					? queueMicrotask
					: Jd === void 0
						? Kd
						: function (e) {
								return Jd.resolve(null).then(e).catch(Xd);
							};
		function Xd(e) {
			setTimeout(function () {
				throw e;
			});
		}
		function Zd(e) {
			return e === `head`;
		}
		function Qd(e, t) {
			var n = t,
				r = 0;
			do {
				var i = n.nextSibling;
				if ((e.removeChild(n), i && i.nodeType === 8))
					if (((n = i.data), n === `/$` || n === `/&`)) {
						if (r === 0) {
							e.removeChild(i), Np(t);
							return;
						}
						r--;
					} else if (n === `$` || n === `$?` || n === `$~` || n === `$!` || n === `&`) r++;
					else if (n === `html`) pf(e.ownerDocument.documentElement);
					else if (n === `head`) {
						(n = e.ownerDocument.head), pf(n);
						for (var a = n.firstChild; a; ) {
							var o = a.nextSibling,
								s = a.nodeName;
							a[xt] ||
								s === `SCRIPT` ||
								s === `STYLE` ||
								(s === `LINK` && a.rel.toLowerCase() === `stylesheet`) ||
								n.removeChild(a),
								(a = o);
						}
					} else n === `body` && pf(e.ownerDocument.body);
				n = i;
			} while (n);
			Np(t);
		}
		function $d(e, t) {
			var n = e;
			e = 0;
			do {
				var r = n.nextSibling;
				if (
					(n.nodeType === 1
						? t
							? ((n._stashedDisplay = n.style.display), (n.style.display = `none`))
							: ((n.style.display = n._stashedDisplay || ``),
								n.getAttribute(`style`) === `` && n.removeAttribute(`style`))
						: n.nodeType === 3 &&
							(t
								? ((n._stashedText = n.nodeValue), (n.nodeValue = ``))
								: (n.nodeValue = n._stashedText || ``)),
					r && r.nodeType === 8)
				)
					if (((n = r.data), n === `/$`)) {
						if (e === 0) break;
						e--;
					} else (n !== `$` && n !== `$?` && n !== `$~` && n !== `$!`) || e++;
				n = r;
			} while (n);
		}
		function ef(e) {
			var t = e.firstChild;
			for (t && t.nodeType === 10 && (t = t.nextSibling); t; ) {
				var n = t;
				switch (((t = t.nextSibling), n.nodeName)) {
					case `HTML`:
					case `HEAD`:
					case `BODY`:
						ef(n), St(n);
						continue;
					case `SCRIPT`:
					case `STYLE`:
						continue;
					case `LINK`:
						if (n.rel.toLowerCase() === `stylesheet`) continue;
				}
				e.removeChild(n);
			}
		}
		function tf(e, t, n, r) {
			for (; e.nodeType === 1; ) {
				var i = n;
				if (e.nodeName.toLowerCase() !== t.toLowerCase()) {
					if (!r && (e.nodeName !== `INPUT` || e.type !== `hidden`)) break;
				} else if (!r)
					if (t === `input` && e.type === `hidden`) {
						var a = i.name == null ? null : `` + i.name;
						if (i.type === `hidden` && e.getAttribute(`name`) === a) return e;
					} else return e;
				else if (!e[xt])
					switch (t) {
						case `meta`:
							if (!e.hasAttribute(`itemprop`)) break;
							return e;
						case `link`:
							if (
								((a = e.getAttribute(`rel`)),
								(a === `stylesheet` && e.hasAttribute(`data-precedence`)) ||
									a !== i.rel ||
									e.getAttribute(`href`) !== (i.href == null || i.href === `` ? null : i.href) ||
									e.getAttribute(`crossorigin`) !==
										(i.crossOrigin == null ? null : i.crossOrigin) ||
									e.getAttribute(`title`) !== (i.title == null ? null : i.title))
							)
								break;
							return e;
						case `style`:
							if (e.hasAttribute(`data-precedence`)) break;
							return e;
						case `script`:
							if (
								((a = e.getAttribute(`src`)),
								(a !== (i.src == null ? null : i.src) ||
									e.getAttribute(`type`) !== (i.type == null ? null : i.type) ||
									e.getAttribute(`crossorigin`) !==
										(i.crossOrigin == null ? null : i.crossOrigin)) &&
									a &&
									e.hasAttribute(`async`) &&
									!e.hasAttribute(`itemprop`))
							)
								break;
							return e;
						default:
							return e;
					}
				if (((e = cf(e.nextSibling)), e === null)) break;
			}
			return null;
		}
		function nf(e, t, n) {
			if (t === ``) return null;
			for (; e.nodeType !== 3; )
				if (
					((e.nodeType !== 1 || e.nodeName !== `INPUT` || e.type !== `hidden`) && !n) ||
					((e = cf(e.nextSibling)), e === null)
				)
					return null;
			return e;
		}
		function rf(e, t) {
			for (; e.nodeType !== 8; )
				if (
					((e.nodeType !== 1 || e.nodeName !== `INPUT` || e.type !== `hidden`) && !t) ||
					((e = cf(e.nextSibling)), e === null)
				)
					return null;
			return e;
		}
		function af(e) {
			return e.data === `$?` || e.data === `$~`;
		}
		function of(e) {
			return e.data === `$!` || (e.data === `$?` && e.ownerDocument.readyState !== `loading`);
		}
		function sf(e, t) {
			var n = e.ownerDocument;
			if (e.data === `$~`) e._reactRetry = t;
			else if (e.data !== `$?` || n.readyState !== `loading`) t();
			else {
				var r = function () {
					t(), n.removeEventListener(`DOMContentLoaded`, r);
				};
				n.addEventListener(`DOMContentLoaded`, r), (e._reactRetry = r);
			}
		}
		function cf(e) {
			for (; e != null; e = e.nextSibling) {
				var t = e.nodeType;
				if (t === 1 || t === 3) break;
				if (t === 8) {
					if (
						((t = e.data),
						t === `$` ||
							t === `$!` ||
							t === `$?` ||
							t === `$~` ||
							t === `&` ||
							t === `F!` ||
							t === `F`)
					)
						break;
					if (t === `/$` || t === `/&`) return null;
				}
			}
			return e;
		}
		var lf = null;
		function uf(e) {
			e = e.nextSibling;
			for (var t = 0; e; ) {
				if (e.nodeType === 8) {
					var n = e.data;
					if (n === `/$` || n === `/&`) {
						if (t === 0) return cf(e.nextSibling);
						t--;
					} else (n !== `$` && n !== `$!` && n !== `$?` && n !== `$~` && n !== `&`) || t++;
				}
				e = e.nextSibling;
			}
			return null;
		}
		function df(e) {
			e = e.previousSibling;
			for (var t = 0; e; ) {
				if (e.nodeType === 8) {
					var n = e.data;
					if (n === `$` || n === `$!` || n === `$?` || n === `$~` || n === `&`) {
						if (t === 0) return e;
						t--;
					} else (n !== `/$` && n !== `/&`) || t++;
				}
				e = e.previousSibling;
			}
			return null;
		}
		function ff(e, t, n) {
			switch (((t = Bd(n)), e)) {
				case `html`:
					if (((e = t.documentElement), !e)) throw Error(s(452));
					return e;
				case `head`:
					if (((e = t.head), !e)) throw Error(s(453));
					return e;
				case `body`:
					if (((e = t.body), !e)) throw Error(s(454));
					return e;
				default:
					throw Error(s(451));
			}
		}
		function pf(e) {
			for (var t = e.attributes; t.length; ) e.removeAttributeNode(t[0]);
			St(e);
		}
		var mf = new Map(),
			hf = new Set();
		function gf(e) {
			return typeof e.getRootNode == `function`
				? e.getRootNode()
				: e.nodeType === 9
					? e
					: e.ownerDocument;
		}
		var _f = E.d;
		E.d = { f: vf, r: yf, D: Sf, C: Cf, L: wf, m: Tf, X: Df, S: Ef, M: Of };
		function vf() {
			var e = _f.f(),
				t = bu();
			return e || t;
		}
		function yf(e) {
			var t = wt(e);
			t !== null && t.tag === 5 && t.type === `form` ? Os(t) : _f.r(e);
		}
		var bf = typeof document > `u` ? null : document;
		function xf(e, t, n) {
			var r = bf;
			if (r && typeof t == `string` && t) {
				var i = Kt(t);
				(i = `link[rel="` + e + `"][href="` + i + `"]`),
					typeof n == `string` && (i += `[crossorigin="` + n + `"]`),
					hf.has(i) ||
						(hf.add(i),
						(e = { rel: e, crossOrigin: n, href: t }),
						r.querySelector(i) === null &&
							((t = r.createElement(`link`)), Pd(t, `link`, e), Dt(t), r.head.appendChild(t)));
			}
		}
		function Sf(e) {
			_f.D(e), xf(`dns-prefetch`, e, null);
		}
		function Cf(e, t) {
			_f.C(e, t), xf(`preconnect`, e, t);
		}
		function wf(e, t, n) {
			_f.L(e, t, n);
			var r = bf;
			if (r && e && t) {
				var i = `link[rel="preload"][as="` + Kt(t) + `"]`;
				t === `image` && n && n.imageSrcSet
					? ((i += `[imagesrcset="` + Kt(n.imageSrcSet) + `"]`),
						typeof n.imageSizes == `string` && (i += `[imagesizes="` + Kt(n.imageSizes) + `"]`))
					: (i += `[href="` + Kt(e) + `"]`);
				var a = i;
				switch (t) {
					case `style`:
						a = Af(e);
						break;
					case `script`:
						a = Pf(e);
				}
				mf.has(a) ||
					((e = h(
						{ rel: `preload`, href: t === `image` && n && n.imageSrcSet ? void 0 : e, as: t },
						n,
					)),
					mf.set(a, e),
					r.querySelector(i) !== null ||
						(t === `style` && r.querySelector(jf(a))) ||
						(t === `script` && r.querySelector(Ff(a))) ||
						((t = r.createElement(`link`)), Pd(t, `link`, e), Dt(t), r.head.appendChild(t)));
			}
		}
		function Tf(e, t) {
			_f.m(e, t);
			var n = bf;
			if (n && e) {
				var r = t && typeof t.as == `string` ? t.as : `script`,
					i = `link[rel="modulepreload"][as="` + Kt(r) + `"][href="` + Kt(e) + `"]`,
					a = i;
				switch (r) {
					case `audioworklet`:
					case `paintworklet`:
					case `serviceworker`:
					case `sharedworker`:
					case `worker`:
					case `script`:
						a = Pf(e);
				}
				if (
					!mf.has(a) &&
					((e = h({ rel: `modulepreload`, href: e }, t)), mf.set(a, e), n.querySelector(i) === null)
				) {
					switch (r) {
						case `audioworklet`:
						case `paintworklet`:
						case `serviceworker`:
						case `sharedworker`:
						case `worker`:
						case `script`:
							if (n.querySelector(Ff(a))) return;
					}
					(r = n.createElement(`link`)), Pd(r, `link`, e), Dt(r), n.head.appendChild(r);
				}
			}
		}
		function Ef(e, t, n) {
			_f.S(e, t, n);
			var r = bf;
			if (r && e) {
				var i = Et(r).hoistableStyles,
					a = Af(e);
				t ||= `default`;
				var o = i.get(a);
				if (!o) {
					var s = { loading: 0, preload: null };
					if ((o = r.querySelector(jf(a)))) s.loading = 5;
					else {
						(e = h({ rel: `stylesheet`, href: e, "data-precedence": t }, n)),
							(n = mf.get(a)) && Rf(e, n);
						var c = (o = r.createElement(`link`));
						Dt(c),
							Pd(c, `link`, e),
							(c._p = new Promise(function (e, t) {
								(c.onload = e), (c.onerror = t);
							})),
							c.addEventListener(`load`, function () {
								s.loading |= 1;
							}),
							c.addEventListener(`error`, function () {
								s.loading |= 2;
							}),
							(s.loading |= 4),
							Lf(o, t, r);
					}
					(o = { type: `stylesheet`, instance: o, count: 1, state: s }), i.set(a, o);
				}
			}
		}
		function Df(e, t) {
			_f.X(e, t);
			var n = bf;
			if (n && e) {
				var r = Et(n).hoistableScripts,
					i = Pf(e),
					a = r.get(i);
				a ||
					((a = n.querySelector(Ff(i))),
					a ||
						((e = h({ src: e, async: !0 }, t)),
						(t = mf.get(i)) && zf(e, t),
						(a = n.createElement(`script`)),
						Dt(a),
						Pd(a, `link`, e),
						n.head.appendChild(a)),
					(a = { type: `script`, instance: a, count: 1, state: null }),
					r.set(i, a));
			}
		}
		function Of(e, t) {
			_f.M(e, t);
			var n = bf;
			if (n && e) {
				var r = Et(n).hoistableScripts,
					i = Pf(e),
					a = r.get(i);
				a ||
					((a = n.querySelector(Ff(i))),
					a ||
						((e = h({ src: e, async: !0, type: `module` }, t)),
						(t = mf.get(i)) && zf(e, t),
						(a = n.createElement(`script`)),
						Dt(a),
						Pd(a, `link`, e),
						n.head.appendChild(a)),
					(a = { type: `script`, instance: a, count: 1, state: null }),
					r.set(i, a));
			}
		}
		function kf(e, t, n, r) {
			var i = (i = ge.current) ? gf(i) : null;
			if (!i) throw Error(s(446));
			switch (e) {
				case `meta`:
				case `title`:
					return null;
				case `style`:
					return typeof n.precedence == `string` && typeof n.href == `string`
						? ((t = Af(n.href)),
							(n = Et(i).hoistableStyles),
							(r = n.get(t)),
							r || ((r = { type: `style`, instance: null, count: 0, state: null }), n.set(t, r)),
							r)
						: { type: `void`, instance: null, count: 0, state: null };
				case `link`:
					if (
						n.rel === `stylesheet` &&
						typeof n.href == `string` &&
						typeof n.precedence == `string`
					) {
						e = Af(n.href);
						var a = Et(i).hoistableStyles,
							o = a.get(e);
						if (
							(o ||
								((i = i.ownerDocument || i),
								(o = {
									type: `stylesheet`,
									instance: null,
									count: 0,
									state: { loading: 0, preload: null },
								}),
								a.set(e, o),
								(a = i.querySelector(jf(e))) && !a._p && ((o.instance = a), (o.state.loading = 5)),
								mf.has(e) ||
									((n = {
										rel: `preload`,
										as: `style`,
										href: n.href,
										crossOrigin: n.crossOrigin,
										integrity: n.integrity,
										media: n.media,
										hrefLang: n.hrefLang,
										referrerPolicy: n.referrerPolicy,
									}),
									mf.set(e, n),
									a || Nf(i, e, n, o.state))),
							t && r === null)
						)
							throw Error(s(528, ``));
						return o;
					}
					if (t && r !== null) throw Error(s(529, ``));
					return null;
				case `script`:
					return (
						(t = n.async),
						(n = n.src),
						typeof n == `string` && t && typeof t != `function` && typeof t != `symbol`
							? ((t = Pf(n)),
								(n = Et(i).hoistableScripts),
								(r = n.get(t)),
								r || ((r = { type: `script`, instance: null, count: 0, state: null }), n.set(t, r)),
								r)
							: { type: `void`, instance: null, count: 0, state: null }
					);
				default:
					throw Error(s(444, e));
			}
		}
		function Af(e) {
			return `href="` + Kt(e) + `"`;
		}
		function jf(e) {
			return `link[rel="stylesheet"][` + e + `]`;
		}
		function Mf(e) {
			return h({}, e, { "data-precedence": e.precedence, precedence: null });
		}
		function Nf(e, t, n, r) {
			e.querySelector(`link[rel="preload"][as="style"][` + t + `]`)
				? (r.loading = 1)
				: ((t = e.createElement(`link`)),
					(r.preload = t),
					t.addEventListener(`load`, function () {
						return (r.loading |= 1);
					}),
					t.addEventListener(`error`, function () {
						return (r.loading |= 2);
					}),
					Pd(t, `link`, n),
					Dt(t),
					e.head.appendChild(t));
		}
		function Pf(e) {
			return `[src="` + Kt(e) + `"]`;
		}
		function Ff(e) {
			return `script[async]` + e;
		}
		function If(e, t, n) {
			if ((t.count++, t.instance === null))
				switch (t.type) {
					case `style`:
						var r = e.querySelector(`style[data-href~="` + Kt(n.href) + `"]`);
						if (r) return (t.instance = r), Dt(r), r;
						var i = h({}, n, {
							"data-href": n.href,
							"data-precedence": n.precedence,
							href: null,
							precedence: null,
						});
						return (
							(r = (e.ownerDocument || e).createElement(`style`)),
							Dt(r),
							Pd(r, `style`, i),
							Lf(r, n.precedence, e),
							(t.instance = r)
						);
					case `stylesheet`:
						i = Af(n.href);
						var a = e.querySelector(jf(i));
						if (a) return (t.state.loading |= 4), (t.instance = a), Dt(a), a;
						(r = Mf(n)),
							(i = mf.get(i)) && Rf(r, i),
							(a = (e.ownerDocument || e).createElement(`link`)),
							Dt(a);
						var o = a;
						return (
							(o._p = new Promise(function (e, t) {
								(o.onload = e), (o.onerror = t);
							})),
							Pd(a, `link`, r),
							(t.state.loading |= 4),
							Lf(a, n.precedence, e),
							(t.instance = a)
						);
					case `script`:
						return (
							(a = Pf(n.src)),
							(i = e.querySelector(Ff(a)))
								? ((t.instance = i), Dt(i), i)
								: ((r = n),
									(i = mf.get(a)) && ((r = h({}, n)), zf(r, i)),
									(e = e.ownerDocument || e),
									(i = e.createElement(`script`)),
									Dt(i),
									Pd(i, `link`, r),
									e.head.appendChild(i),
									(t.instance = i))
						);
					case `void`:
						return null;
					default:
						throw Error(s(443, t.type));
				}
			else
				t.type === `stylesheet` &&
					!(t.state.loading & 4) &&
					((r = t.instance), (t.state.loading |= 4), Lf(r, n.precedence, e));
			return t.instance;
		}
		function Lf(e, t, n) {
			for (
				var r = n.querySelectorAll(
						`link[rel="stylesheet"][data-precedence],style[data-precedence]`,
					),
					i = r.length ? r[r.length - 1] : null,
					a = i,
					o = 0;
				o < r.length;
				o++
			) {
				var s = r[o];
				if (s.dataset.precedence === t) a = s;
				else if (a !== i) break;
			}
			a
				? a.parentNode.insertBefore(e, a.nextSibling)
				: ((t = n.nodeType === 9 ? n.head : n), t.insertBefore(e, t.firstChild));
		}
		function Rf(e, t) {
			(e.crossOrigin ??= t.crossOrigin),
				(e.referrerPolicy ??= t.referrerPolicy),
				(e.title ??= t.title);
		}
		function zf(e, t) {
			(e.crossOrigin ??= t.crossOrigin),
				(e.referrerPolicy ??= t.referrerPolicy),
				(e.integrity ??= t.integrity);
		}
		var Bf = null;
		function Vf(e, t, n) {
			if (Bf === null) {
				var r = new Map(),
					i = (Bf = new Map());
				i.set(n, r);
			} else (i = Bf), (r = i.get(n)), r || ((r = new Map()), i.set(n, r));
			if (r.has(e)) return r;
			for (r.set(e, null), n = n.getElementsByTagName(e), i = 0; i < n.length; i++) {
				var a = n[i];
				if (
					!(a[xt] || a[mt] || (e === `link` && a.getAttribute(`rel`) === `stylesheet`)) &&
					a.namespaceURI !== `http://www.w3.org/2000/svg`
				) {
					var o = a.getAttribute(t) || ``;
					o = e + o;
					var s = r.get(o);
					s ? s.push(a) : r.set(o, [a]);
				}
			}
			return r;
		}
		function Hf(e, t, n) {
			(e = e.ownerDocument || e),
				e.head.insertBefore(n, t === `title` ? e.querySelector(`head > title`) : null);
		}
		function Uf(e, t, n) {
			if (n === 1 || t.itemProp != null) return !1;
			switch (e) {
				case `meta`:
				case `title`:
					return !0;
				case `style`:
					if (typeof t.precedence != `string` || typeof t.href != `string` || t.href === ``) break;
					return !0;
				case `link`:
					if (
						typeof t.rel != `string` ||
						typeof t.href != `string` ||
						t.href === `` ||
						t.onLoad ||
						t.onError
					)
						break;
					switch (t.rel) {
						case `stylesheet`:
							return (e = t.disabled), typeof t.precedence == `string` && e == null;
						default:
							return !0;
					}
				case `script`:
					if (
						t.async &&
						typeof t.async != `function` &&
						typeof t.async != `symbol` &&
						!t.onLoad &&
						!t.onError &&
						t.src &&
						typeof t.src == `string`
					)
						return !0;
			}
			return !1;
		}
		function Wf(e) {
			return !(e.type === `stylesheet` && !(e.state.loading & 3));
		}
		function Gf(e, t, n, r) {
			if (
				n.type === `stylesheet` &&
				(typeof r.media != `string` || !1 !== matchMedia(r.media).matches) &&
				!(n.state.loading & 4)
			) {
				if (n.instance === null) {
					var i = Af(r.href),
						a = t.querySelector(jf(i));
					if (a) {
						(t = a._p),
							typeof t == `object` &&
								t &&
								typeof t.then == `function` &&
								(e.count++, (e = Jf.bind(e)), t.then(e, e)),
							(n.state.loading |= 4),
							(n.instance = a),
							Dt(a);
						return;
					}
					(a = t.ownerDocument || t),
						(r = Mf(r)),
						(i = mf.get(i)) && Rf(r, i),
						(a = a.createElement(`link`)),
						Dt(a);
					var o = a;
					(o._p = new Promise(function (e, t) {
						(o.onload = e), (o.onerror = t);
					})),
						Pd(a, `link`, r),
						(n.instance = a);
				}
				e.stylesheets === null && (e.stylesheets = new Map()),
					e.stylesheets.set(n, t),
					(t = n.state.preload) &&
						!(n.state.loading & 3) &&
						(e.count++,
						(n = Jf.bind(e)),
						t.addEventListener(`load`, n),
						t.addEventListener(`error`, n));
			}
		}
		var Kf = 0;
		function qf(e, t) {
			return (
				e.stylesheets && e.count === 0 && Xf(e, e.stylesheets),
				0 < e.count || 0 < e.imgCount
					? function (n) {
							var r = setTimeout(function () {
								if ((e.stylesheets && Xf(e, e.stylesheets), e.unsuspend)) {
									var t = e.unsuspend;
									(e.unsuspend = null), t();
								}
							}, 6e4 + t);
							0 < e.imgBytes && Kf === 0 && (Kf = 62500 * Ld());
							var i = setTimeout(
								function () {
									if (
										((e.waitingForImages = !1),
										e.count === 0 && (e.stylesheets && Xf(e, e.stylesheets), e.unsuspend))
									) {
										var t = e.unsuspend;
										(e.unsuspend = null), t();
									}
								},
								(e.imgBytes > Kf ? 50 : 800) + t,
							);
							return (
								(e.unsuspend = n),
								function () {
									(e.unsuspend = null), clearTimeout(r), clearTimeout(i);
								}
							);
						}
					: null
			);
		}
		function Jf() {
			if ((this.count--, this.count === 0 && (this.imgCount === 0 || !this.waitingForImages))) {
				if (this.stylesheets) Xf(this, this.stylesheets);
				else if (this.unsuspend) {
					var e = this.unsuspend;
					(this.unsuspend = null), e();
				}
			}
		}
		var Yf = null;
		function Xf(e, t) {
			(e.stylesheets = null),
				e.unsuspend !== null &&
					(e.count++, (Yf = new Map()), t.forEach(Zf, e), (Yf = null), Jf.call(e));
		}
		function Zf(e, t) {
			if (!(t.state.loading & 4)) {
				var n = Yf.get(e);
				if (n) var r = n.get(null);
				else {
					(n = new Map()), Yf.set(e, n);
					for (
						var i = e.querySelectorAll(`link[data-precedence],style[data-precedence]`), a = 0;
						a < i.length;
						a++
					) {
						var o = i[a];
						(o.nodeName === `LINK` || o.getAttribute(`media`) !== `not all`) &&
							(n.set(o.dataset.precedence, o), (r = o));
					}
					r && n.set(null, r);
				}
				(i = t.instance),
					(o = i.getAttribute(`data-precedence`)),
					(a = n.get(o) || r),
					a === r && n.set(null, i),
					n.set(o, i),
					this.count++,
					(r = Jf.bind(this)),
					i.addEventListener(`load`, r),
					i.addEventListener(`error`, r),
					a
						? a.parentNode.insertBefore(i, a.nextSibling)
						: ((e = e.nodeType === 9 ? e.head : e), e.insertBefore(i, e.firstChild)),
					(t.state.loading |= 4);
			}
		}
		var Qf = {
			$$typeof: S,
			Provider: null,
			Consumer: null,
			_currentValue: de,
			_currentValue2: de,
			_threadCount: 0,
		};
		function $f(e, t, n, r, i, a, o, s, c) {
			(this.tag = 1),
				(this.containerInfo = e),
				(this.pingCache = this.current = this.pendingChildren = null),
				(this.timeoutHandle = -1),
				(this.callbackNode =
					this.next =
					this.pendingContext =
					this.context =
					this.cancelPendingCommit =
						null),
				(this.callbackPriority = 0),
				(this.expirationTimes = rt(-1)),
				(this.entangledLanes =
					this.shellSuspendCounter =
					this.errorRecoveryDisabledLanes =
					this.expiredLanes =
					this.warmLanes =
					this.pingedLanes =
					this.suspendedLanes =
					this.pendingLanes =
						0),
				(this.entanglements = rt(0)),
				(this.hiddenUpdates = rt(null)),
				(this.identifierPrefix = r),
				(this.onUncaughtError = i),
				(this.onCaughtError = a),
				(this.onRecoverableError = o),
				(this.pooledCache = null),
				(this.pooledCacheLanes = 0),
				(this.formState = c),
				(this.incompleteTransitions = new Map());
		}
		function ep(e, t, n, r, i, a, o, s, c, l, u, d) {
			return (
				(e = new $f(e, t, n, o, c, l, u, d, s)),
				(t = 1),
				!0 === a && (t |= 24),
				(a = mi(3, null, null, t)),
				(e.current = a),
				(a.stateNode = e),
				(t = fa()),
				t.refCount++,
				(e.pooledCache = t),
				t.refCount++,
				(a.memoizedState = { element: r, isDehydrated: n, cache: t }),
				Ga(a),
				e
			);
		}
		function tp(e) {
			return e ? ((e = fi), e) : fi;
		}
		function np(e, t, n, r, i, a) {
			(i = tp(i)),
				r.context === null ? (r.context = i) : (r.pendingContext = i),
				(r = qa(t)),
				(r.payload = { element: n }),
				(a = a === void 0 ? null : a),
				a !== null && (r.callback = a),
				(n = Ja(e, r, t)),
				n !== null && (hu(n, e, t), Ya(n, e, t));
		}
		function rp(e, t) {
			if (((e = e.memoizedState), e !== null && e.dehydrated !== null)) {
				var n = e.retryLane;
				e.retryLane = n !== 0 && n < t ? n : t;
			}
		}
		function ip(e, t) {
			rp(e, t), (e = e.alternate) && rp(e, t);
		}
		function ap(e) {
			if (e.tag === 13 || e.tag === 31) {
				var t = li(e, 67108864);
				t !== null && hu(t, e, 67108864), ip(e, 67108864);
			}
		}
		function op(e) {
			if (e.tag === 13 || e.tag === 31) {
				var t = pu();
				t = lt(t);
				var n = li(e, t);
				n !== null && hu(n, e, t), ip(e, t);
			}
		}
		var sp = !0;
		function cp(e, t, n, r) {
			var i = T.T;
			T.T = null;
			var a = E.p;
			try {
				(E.p = 2), up(e, t, n, r);
			} finally {
				(E.p = a), (T.T = i);
			}
		}
		function lp(e, t, n, r) {
			var i = T.T;
			T.T = null;
			var a = E.p;
			try {
				(E.p = 8), up(e, t, n, r);
			} finally {
				(E.p = a), (T.T = i);
			}
		}
		function up(e, t, n, r) {
			if (sp) {
				var i = dp(r);
				if (i === null) wd(e, t, r, fp, n), Cp(e, r);
				else if (Tp(i, e, t, n, r)) r.stopPropagation();
				else if ((Cp(e, r), t & 4 && -1 < Sp.indexOf(e))) {
					for (; i !== null; ) {
						var a = wt(i);
						if (a !== null)
							switch (a.tag) {
								case 3:
									if (((a = a.stateNode), a.current.memoizedState.isDehydrated)) {
										var o = Qe(a.pendingLanes);
										if (o !== 0) {
											var s = a;
											for (s.pendingLanes |= 2, s.entangledLanes |= 2; o; ) {
												var c = 1 << (31 - Ge(o));
												(s.entanglements[1] |= c), (o &= ~c);
											}
											rd(a), !(G & 6) && ((tu = Ne() + 500), id(0, !1));
										}
									}
									break;
								case 31:
								case 13:
									(s = li(a, 2)), s !== null && hu(s, a, 2), bu(), ip(a, 2);
							}
						if (((a = dp(r)), a === null && wd(e, t, r, fp, n), a === i)) break;
						i = a;
					}
					i !== null && r.stopPropagation();
				} else wd(e, t, r, null, n);
			}
		}
		function dp(e) {
			return (e = ln(e)), pp(e);
		}
		var fp = null;
		function pp(e) {
			if (((fp = null), (e = Ct(e)), e !== null)) {
				var t = l(e);
				if (t === null) e = null;
				else {
					var n = t.tag;
					if (n === 13) {
						if (((e = u(t)), e !== null)) return e;
						e = null;
					} else if (n === 31) {
						if (((e = d(t)), e !== null)) return e;
						e = null;
					} else if (n === 3) {
						if (t.stateNode.current.memoizedState.isDehydrated)
							return t.tag === 3 ? t.stateNode.containerInfo : null;
						e = null;
					} else t !== e && (e = null);
				}
			}
			return (fp = e), null;
		}
		function mp(e) {
			switch (e) {
				case `beforetoggle`:
				case `cancel`:
				case `click`:
				case `close`:
				case `contextmenu`:
				case `copy`:
				case `cut`:
				case `auxclick`:
				case `dblclick`:
				case `dragend`:
				case `dragstart`:
				case `drop`:
				case `focusin`:
				case `focusout`:
				case `input`:
				case `invalid`:
				case `keydown`:
				case `keypress`:
				case `keyup`:
				case `mousedown`:
				case `mouseup`:
				case `paste`:
				case `pause`:
				case `play`:
				case `pointercancel`:
				case `pointerdown`:
				case `pointerup`:
				case `ratechange`:
				case `reset`:
				case `resize`:
				case `seeked`:
				case `submit`:
				case `toggle`:
				case `touchcancel`:
				case `touchend`:
				case `touchstart`:
				case `volumechange`:
				case `change`:
				case `selectionchange`:
				case `textInput`:
				case `compositionstart`:
				case `compositionend`:
				case `compositionupdate`:
				case `beforeblur`:
				case `afterblur`:
				case `beforeinput`:
				case `blur`:
				case `fullscreenchange`:
				case `focus`:
				case `hashchange`:
				case `popstate`:
				case `select`:
				case `selectstart`:
					return 2;
				case `drag`:
				case `dragenter`:
				case `dragexit`:
				case `dragleave`:
				case `dragover`:
				case `mousemove`:
				case `mouseout`:
				case `mouseover`:
				case `pointermove`:
				case `pointerout`:
				case `pointerover`:
				case `scroll`:
				case `touchmove`:
				case `wheel`:
				case `mouseenter`:
				case `mouseleave`:
				case `pointerenter`:
				case `pointerleave`:
					return 8;
				case `message`:
					switch (Pe()) {
						case Fe:
							return 2;
						case Ie:
							return 8;
						case Le:
						case Re:
							return 32;
						case ze:
							return 268435456;
						default:
							return 32;
					}
				default:
					return 32;
			}
		}
		var hp = !1,
			gp = null,
			_p = null,
			vp = null,
			yp = new Map(),
			bp = new Map(),
			xp = [],
			Sp =
				`mousedown mouseup touchcancel touchend touchstart auxclick dblclick pointercancel pointerdown pointerup dragend dragstart drop compositionend compositionstart keydown keypress keyup input textInput copy cut paste click change contextmenu reset`.split(
					` `,
				);
		function Cp(e, t) {
			switch (e) {
				case `focusin`:
				case `focusout`:
					gp = null;
					break;
				case `dragenter`:
				case `dragleave`:
					_p = null;
					break;
				case `mouseover`:
				case `mouseout`:
					vp = null;
					break;
				case `pointerover`:
				case `pointerout`:
					yp.delete(t.pointerId);
					break;
				case `gotpointercapture`:
				case `lostpointercapture`:
					bp.delete(t.pointerId);
			}
		}
		function wp(e, t, n, r, i, a) {
			return e === null || e.nativeEvent !== a
				? ((e = {
						blockedOn: t,
						domEventName: n,
						eventSystemFlags: r,
						nativeEvent: a,
						targetContainers: [i],
					}),
					t !== null && ((t = wt(t)), t !== null && ap(t)),
					e)
				: ((e.eventSystemFlags |= r),
					(t = e.targetContainers),
					i !== null && t.indexOf(i) === -1 && t.push(i),
					e);
		}
		function Tp(e, t, n, r, i) {
			switch (t) {
				case `focusin`:
					return (gp = wp(gp, e, t, n, r, i)), !0;
				case `dragenter`:
					return (_p = wp(_p, e, t, n, r, i)), !0;
				case `mouseover`:
					return (vp = wp(vp, e, t, n, r, i)), !0;
				case `pointerover`:
					var a = i.pointerId;
					return yp.set(a, wp(yp.get(a) || null, e, t, n, r, i)), !0;
				case `gotpointercapture`:
					return (a = i.pointerId), bp.set(a, wp(bp.get(a) || null, e, t, n, r, i)), !0;
			}
			return !1;
		}
		function Ep(e) {
			var t = Ct(e.target);
			if (t !== null) {
				var n = l(t);
				if (n !== null) {
					if (((t = n.tag), t === 13)) {
						if (((t = u(n)), t !== null)) {
							(e.blockedOn = t),
								ft(e.priority, function () {
									op(n);
								});
							return;
						}
					} else if (t === 31) {
						if (((t = d(n)), t !== null)) {
							(e.blockedOn = t),
								ft(e.priority, function () {
									op(n);
								});
							return;
						}
					} else if (t === 3 && n.stateNode.current.memoizedState.isDehydrated) {
						e.blockedOn = n.tag === 3 ? n.stateNode.containerInfo : null;
						return;
					}
				}
			}
			e.blockedOn = null;
		}
		function Dp(e) {
			if (e.blockedOn !== null) return !1;
			for (var t = e.targetContainers; 0 < t.length; ) {
				var n = dp(e.nativeEvent);
				if (n === null) {
					n = e.nativeEvent;
					var r = new n.constructor(n.type, n);
					(cn = r), n.target.dispatchEvent(r), (cn = null);
				} else return (t = wt(n)), t !== null && ap(t), (e.blockedOn = n), !1;
				t.shift();
			}
			return !0;
		}
		function Op(e, t, n) {
			Dp(e) && n.delete(t);
		}
		function kp() {
			(hp = !1),
				gp !== null && Dp(gp) && (gp = null),
				_p !== null && Dp(_p) && (_p = null),
				vp !== null && Dp(vp) && (vp = null),
				yp.forEach(Op),
				bp.forEach(Op);
		}
		function Ap(e, n) {
			e.blockedOn === n &&
				((e.blockedOn = null),
				hp || ((hp = !0), t.unstable_scheduleCallback(t.unstable_NormalPriority, kp)));
		}
		var jp = null;
		function Mp(e) {
			jp !== e &&
				((jp = e),
				t.unstable_scheduleCallback(t.unstable_NormalPriority, function () {
					jp === e && (jp = null);
					for (var t = 0; t < e.length; t += 3) {
						var n = e[t],
							r = e[t + 1],
							i = e[t + 2];
						if (typeof r != `function`) {
							if (pp(r || n) === null) continue;
							break;
						}
						var a = wt(n);
						a !== null &&
							(e.splice(t, 3),
							(t -= 3),
							Es(a, { pending: !0, data: i, method: n.method, action: r }, r, i));
					}
				}));
		}
		function Np(e) {
			function t(t) {
				return Ap(t, e);
			}
			gp !== null && Ap(gp, e),
				_p !== null && Ap(_p, e),
				vp !== null && Ap(vp, e),
				yp.forEach(t),
				bp.forEach(t);
			for (var n = 0; n < xp.length; n++) {
				var r = xp[n];
				r.blockedOn === e && (r.blockedOn = null);
			}
			for (; 0 < xp.length && ((n = xp[0]), n.blockedOn === null); )
				Ep(n), n.blockedOn === null && xp.shift();
			if (((n = (e.ownerDocument || e).$$reactFormReplay), n != null))
				for (r = 0; r < n.length; r += 3) {
					var i = n[r],
						a = n[r + 1],
						o = i[ht] || null;
					if (typeof a == `function`) o || Mp(n);
					else if (o) {
						var s = null;
						if (a && a.hasAttribute(`formAction`)) {
							if (((i = a), (o = a[ht] || null))) s = o.formAction;
							else if (pp(i) !== null) continue;
						} else s = o.action;
						typeof s == `function` ? (n[r + 1] = s) : (n.splice(r, 3), (r -= 3)), Mp(n);
					}
				}
		}
		function Pp() {
			function e(e) {
				e.canIntercept &&
					e.info === `react-transition` &&
					e.intercept({
						handler: function () {
							return new Promise(function (e) {
								return (i = e);
							});
						},
						focusReset: `manual`,
						scroll: `manual`,
					});
			}
			function t() {
				i !== null && (i(), (i = null)), r || setTimeout(n, 20);
			}
			function n() {
				if (!r && !navigation.transition) {
					var e = navigation.currentEntry;
					e &&
						e.url != null &&
						navigation.navigate(e.url, {
							state: e.getState(),
							info: `react-transition`,
							history: `replace`,
						});
				}
			}
			if (typeof navigation == `object`) {
				var r = !1,
					i = null;
				return (
					navigation.addEventListener(`navigate`, e),
					navigation.addEventListener(`navigatesuccess`, t),
					navigation.addEventListener(`navigateerror`, t),
					setTimeout(n, 100),
					function () {
						(r = !0),
							navigation.removeEventListener(`navigate`, e),
							navigation.removeEventListener(`navigatesuccess`, t),
							navigation.removeEventListener(`navigateerror`, t),
							i !== null && (i(), (i = null));
					}
				);
			}
		}
		function Fp(e) {
			this._internalRoot = e;
		}
		(Ip.prototype.render = Fp.prototype.render =
			function (e) {
				var t = this._internalRoot;
				if (t === null) throw Error(s(409));
				var n = t.current;
				np(n, pu(), e, t, null, null);
			}),
			(Ip.prototype.unmount = Fp.prototype.unmount =
				function () {
					var e = this._internalRoot;
					if (e !== null) {
						this._internalRoot = null;
						var t = e.containerInfo;
						np(e.current, 2, null, e, null, null), bu(), (t[gt] = null);
					}
				});
		function Ip(e) {
			this._internalRoot = e;
		}
		Ip.prototype.unstable_scheduleHydration = function (e) {
			if (e) {
				var t = dt();
				e = { blockedOn: null, target: e, priority: t };
				for (var n = 0; n < xp.length && t !== 0 && t < xp[n].priority; n++);
				xp.splice(n, 0, e), n === 0 && Ep(e);
			}
		};
		var Lp = r.version;
		if (Lp !== `19.2.7`) throw Error(s(527, Lp, `19.2.7`));
		E.findDOMNode = function (e) {
			var t = e._reactInternals;
			if (t === void 0)
				throw typeof e.render == `function`
					? Error(s(188))
					: ((e = Object.keys(e).join(`,`)), Error(s(268, e)));
			return (e = p(t)), (e = e === null ? null : m(e)), (e = e === null ? null : e.stateNode), e;
		};
		var Rp = {
			bundleType: 0,
			version: `19.2.7`,
			rendererPackageName: `react-dom`,
			currentDispatcherRef: T,
			reconcilerVersion: `19.2.7`,
		};
		if (typeof __REACT_DEVTOOLS_GLOBAL_HOOK__ < `u`) {
			var zp = __REACT_DEVTOOLS_GLOBAL_HOOK__;
			if (!zp.isDisabled && zp.supportsFiber)
				try {
					(He = zp.inject(Rp)), (Ue = zp);
				} catch {}
		}
		e.createRoot = function (e, t) {
			if (!c(e)) throw Error(s(299));
			var n = !1,
				r = ``,
				i = Ys,
				a = Xs,
				o = Zs;
			return (
				t != null &&
					(!0 === t.unstable_strictMode && (n = !0),
					t.identifierPrefix !== void 0 && (r = t.identifierPrefix),
					t.onUncaughtError !== void 0 && (i = t.onUncaughtError),
					t.onCaughtError !== void 0 && (a = t.onCaughtError),
					t.onRecoverableError !== void 0 && (o = t.onRecoverableError)),
				(t = ep(e, 1, !1, null, null, n, r, null, i, a, o, Pp)),
				(e[gt] = t.current),
				Sd(e),
				new Fp(t)
			);
		};
	}),
	c = e((e, t) => {
		function n() {
			if (
				!(
					typeof __REACT_DEVTOOLS_GLOBAL_HOOK__ > `u` ||
					typeof __REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE != `function`
				)
			)
				try {
					__REACT_DEVTOOLS_GLOBAL_HOOK__.checkDCE(n);
				} catch (e) {
					console.error(e);
				}
		}
		n(), (t.exports = s());
	}),
	l = {
		root: `_root_1gwp3_1`,
		primary: `_primary_1gwp3_20`,
		selected: `_selected_1gwp3_21`,
		compact: `_compact_1gwp3_36`,
		normal: `_normal_1gwp3_41`,
	},
	u = e((e) => {
		var t = Symbol.for(`react.transitional.element`),
			n = Symbol.for(`react.fragment`);
		function r(e, n, r) {
			var i = null;
			if ((r !== void 0 && (i = `` + r), n.key !== void 0 && (i = `` + n.key), `key` in n))
				for (var a in ((r = {}), n)) a !== `key` && (r[a] = n[a]);
			else r = n;
			return (n = r.ref), { $$typeof: t, type: e, key: i, ref: n === void 0 ? null : n, props: r };
		}
		(e.Fragment = n), (e.jsx = r), (e.jsxs = r);
	}),
	d = e((e, t) => {
		t.exports = u();
	})();
function f({
	variant: e = `default`,
	size: t = `normal`,
	selected: n = !1,
	className: r,
	children: i,
	...a
}) {
	return (0, d.jsx)(`button`, {
		type: `button`,
		className: [
			l.root,
			e === `primary` ? l.primary : ``,
			n ? l.selected : ``,
			t === `compact` ? l.compact : l.normal,
			r ?? ``,
		]
			.filter(Boolean)
			.join(` `),
		"data-rag-atom": `Button`,
		...a,
		children: i,
	});
}
var p = { root: `_root_e3xd9_1` };
function m({ className: e, children: t, ...n }) {
	return (0, d.jsx)(`div`, {
		className: [p.root, e ?? ``].filter(Boolean).join(` `),
		role: `alert`,
		"data-rag-atom": `ErrorCallout`,
		...n,
		children: t,
	});
}
var h = { root: `_root_im5yi_1` };
function g({ className: e, children: t, ...n }) {
	return (0, d.jsx)(`select`, {
		className: [h.root, e ?? ``].filter(Boolean).join(` `),
		"data-rag-atom": `SelectInput`,
		...n,
		children: t,
	});
}
var _ = { root: `_root_1bc5y_1` };
function v({ className: e, ...t }) {
	return (0, d.jsx)(`input`, {
		className: [_.root, e ?? ``].filter(Boolean).join(` `),
		"data-rag-atom": `TextInput`,
		...t,
	});
}
var y = `#000000`,
	b = `#0000CC`,
	x = [
		`system`,
		`instruction`,
		`context`,
		`conversation`,
		`summary`,
		`retrieval`,
		`tool`,
		`result`,
		`generated`,
		`annotation`,
		`course`,
		`evicted`,
		`active`,
		`empty`,
		`other`,
	],
	ee = {
		empty: { fill: `#ffffff`, stroke: y, dashed: !0, label: `#666666`, name: `free space` },
		context: { fill: `#ffffff`, stroke: y, label: y, name: `context / history` },
		conversation: { fill: `#ffffff`, stroke: y, label: y, name: `conversation` },
		instruction: {
			fill: `var(--rag-context-pattern-checker)`,
			stroke: y,
			label: y,
			name: `instruction`,
			pattern: `checker`,
		},
		system: {
			fill: `var(--rag-context-pattern-checker)`,
			stroke: y,
			label: y,
			name: `system / instructions`,
			pattern: `checker`,
		},
		summary: {
			fill: `var(--rag-context-pattern-diagonal)`,
			stroke: y,
			label: y,
			name: `summary / memory`,
			pattern: `diagonal`,
		},
		retrieval: {
			fill: `var(--rag-context-pattern-stipple)`,
			stroke: y,
			label: y,
			name: `retrieval`,
			pattern: `stipple`,
		},
		tool: {
			fill: `var(--rag-context-pattern-stipple)`,
			stroke: y,
			label: y,
			name: `tool call`,
			pattern: `stipple`,
		},
		result: {
			fill: `var(--rag-context-pattern-stipple)`,
			stroke: y,
			label: y,
			name: `tool result`,
			pattern: `stipple`,
		},
		generated: {
			fill: `var(--rag-context-pattern-diagonal-dense)`,
			stroke: y,
			label: y,
			name: `generated / scratchpad`,
			pattern: `diagonalDense`,
		},
		annotation: {
			fill: `var(--rag-context-pattern-diagonal)`,
			stroke: y,
			label: y,
			name: `annotation`,
			pattern: `diagonal`,
		},
		course: { fill: `#ffffff`, stroke: y, label: y, name: `course material` },
		evicted: {
			fill: `var(--rag-context-pattern-cross)`,
			stroke: y,
			label: y,
			name: `evicted / compressed`,
			pattern: `cross`,
		},
		active: { fill: b, stroke: y, label: `#ffffff`, name: `active (current)` },
		other: { fill: `#ffffff`, stroke: y, label: y, name: `other` },
	},
	S = {
		empty: { fill: `#ffffff`, dashed: !0, label: `#666666` },
		context: { fill: `#ffffff`, label: `#000000` },
		conversation: { fill: `#ffffff`, label: `#000000` },
		instruction: { fill: `#dddddd`, label: `#000000`, pattern: void 0 },
		system: { fill: `#dddddd`, label: `#000000`, pattern: void 0 },
		summary: { fill: `#eeeeee`, label: `#000000`, pattern: void 0 },
		retrieval: { fill: `#d2d2d2`, label: `#000000`, pattern: void 0 },
		tool: { fill: `#d2d2d2`, label: `#000000`, pattern: void 0 },
		result: { fill: `#c4c4c4`, label: `#000000`, pattern: void 0 },
		generated: { fill: `#a0a0a0`, label: `#000000`, pattern: void 0 },
		annotation: { fill: `#eeeeee`, label: `#000000`, pattern: void 0 },
		course: { fill: `#ffffff`, label: `#000000` },
		evicted: { fill: `#6a6a6a`, label: `#ffffff`, pattern: void 0 },
		active: { fill: b, label: `#ffffff` },
		other: { fill: `#ffffff`, label: `#000000` },
	},
	C = {
		empty: { fill: `#ffffff`, stroke: `#999999`, dashed: !0, label: `#666666` },
		context: { fill: `#ffffff`, strokeWidth: 1.25, label: `#000000`, pattern: void 0 },
		conversation: { fill: `#ffffff`, strokeWidth: 1.25, label: `#000000`, pattern: void 0 },
		instruction: { fill: `#ffffff`, strokeWidth: 2.5, label: `#000000`, pattern: void 0 },
		system: { fill: `#ffffff`, strokeWidth: 2.5, label: `#000000`, pattern: void 0 },
		summary: { fill: `#ffffff`, dashed: !0, label: `#000000`, pattern: void 0 },
		retrieval: { fill: `#ffffff`, dotted: !0, label: `#000000`, pattern: void 0 },
		tool: { fill: `#ffffff`, dotted: !0, label: `#000000`, pattern: void 0 },
		result: { fill: `#ffffff`, dotted: !0, label: `#000000`, pattern: void 0 },
		generated: { fill: `#ffffff`, strokeWidth: 2.5, dashed: !0, label: `#000000`, pattern: void 0 },
		annotation: { fill: `#ffffff`, dashed: !0, label: `#000000`, pattern: void 0 },
		course: { fill: `#ffffff`, strokeWidth: 1.25, label: `#000000` },
		evicted: { fill: `#ffffff`, inner: !0, dashed: !0, label: `#000000`, pattern: void 0 },
		active: { fill: b, label: `#ffffff` },
		other: { fill: `#ffffff`, label: `#000000` },
	};
function te(e, t = `pattern`) {
	let n = re(e) ? e : `other`,
		r = ee[n];
	return t === `tone`
		? { ...r, ...S[n], stroke: S[n]?.stroke ?? `#000000` }
		: t === `outline`
			? { ...r, ...C[n], stroke: C[n]?.stroke ?? `#000000` }
			: r;
}
function ne(e) {
	return te(e).name;
}
function re(e) {
	return x.includes(e);
}
function w(e, t = {}) {
	return e.parts.reduce((e, n) => (!t.includeEmpty && n.kind === `empty` ? e : e + n.tokens), 0);
}
function ie(e) {
	return e.limit <= 0 ? 0 : w(e) / e.limit;
}
var ae = {
	root: `_root_7tnun_1`,
	sm: `_sm_7tnun_11`,
	md: `_md_7tnun_16`,
	lg: `_lg_7tnun_21`,
	selected: `_selected_7tnun_26`,
	dashed: `_dashed_7tnun_31`,
	dotted: `_dotted_7tnun_35`,
	inner: `_inner_7tnun_39`,
	checker: `_checker_7tnun_43`,
	diagonal: `_diagonal_7tnun_54`,
	diagonalDense: `_diagonalDense_7tnun_59`,
	stipple: `_stipple_7tnun_64`,
	cross: `_cross_7tnun_70`,
	overflow: `_overflow_7tnun_77`,
};
function oe({
	kind: e,
	mode: t = `pattern`,
	size: n = `md`,
	selected: r = !1,
	className: i,
	style: a,
	...o
}) {
	let s = te(e, t),
		c = ne(e);
	return (0, d.jsx)(`span`, {
		className: [
			ae.root,
			ae[n],
			s.pattern ? ae[s.pattern] : ``,
			s.dashed ? ae.dashed : ``,
			s.dotted ? ae.dotted : ``,
			s.inner ? ae.inner : ``,
			r ? ae.selected : ``,
			i ?? ``,
		]
			.filter(Boolean)
			.join(` `),
		"data-rag-atom": `ContextKindSwatch`,
		"data-kind": e,
		"data-mode": t,
		"aria-label": c,
		title: c,
		style: {
			"--rag-context-swatch-fill": s.fill,
			"--rag-context-swatch-stroke": s.stroke,
			"--rag-context-swatch-stroke-width": `${s.strokeWidth ?? 1}px`,
			...a,
		},
		...o,
	});
}
var se = {
	root: `_root_j1d2u_1`,
	user: `_user_j1d2u_2`,
	assistant: `_assistant_j1d2u_3`,
	tool: `_tool_j1d2u_4`,
	system: `_system_j1d2u_5`,
	developer: `_developer_j1d2u_5`,
	other: `_other_j1d2u_6`,
};
function ce({ role: e, name: t, className: n, ...r }) {
	return (0, d.jsx)(`span`, {
		className: [se.root, se[e] ?? se.other, n ?? ``].filter(Boolean).join(` `),
		"data-rag-atom": `TranscriptRoleBadge`,
		"data-role": e,
		...r,
		children: t ?? e,
	});
}
var le = { root: `_root_1dvu0_1`, selected: `_selected_1dvu0_2` };
function ue({ kind: e, label: t, selected: n = !1, className: r, ...i }) {
	return (0, d.jsxs)(`span`, {
		className: [le.root, n ? le.selected : ``, r ?? ``].filter(Boolean).join(` `),
		"data-rag-atom": `AnnotationBadge`,
		"data-kind": e,
		...i,
		children: [
			(0, d.jsx)(oe, { kind: e, size: `sm`, selected: n }),
			(0, d.jsx)(`span`, { children: t }),
		],
	});
}
var T = { root: `_root_72kwi_1`, folder: `_folder_72kwi_8`, arrow: `_arrow_72kwi_15` };
function E({ title: e = `Upload file`, className: t, ...n }) {
	let r = e ? `rag-upload-glyph-title` : void 0;
	return (0, d.jsxs)(`svg`, {
		className: [T.root, t ?? ``].filter(Boolean).join(` `),
		viewBox: `0 0 40 34`,
		role: e ? `img` : `presentation`,
		"aria-labelledby": r,
		"data-rag-atom": `UploadGlyph`,
		...n,
		children: [
			e && (0, d.jsx)(`title`, { id: r, children: e }),
			(0, d.jsx)(`path`, { className: T.folder, d: `M2 8 L2 32 L38 32 L38 4 L18 4 L14 8 Z` }),
			(0, d.jsx)(`line`, { className: T.arrow, x1: `20`, y1: `14`, x2: `20`, y2: `26` }),
			(0, d.jsx)(`path`, { className: T.arrow, d: `M15 19 L20 14 L25 19` }),
		],
	});
}
var de = { root: `_root_1gd53_1`, filled: `_filled_1gd53_13` };
function fe({ id: e, title: t, className: n, ...r }) {
	let i = t ? `rag-context-studio-nav-icon-${e}` : void 0;
	return (0, d.jsxs)(`svg`, {
		className: [de.root, n ?? ``].filter(Boolean).join(` `),
		viewBox: `0 0 15 15`,
		role: t ? `img` : `presentation`,
		"aria-labelledby": i,
		"data-rag-atom": `ContextStudioNavIcon`,
		"data-icon-id": e,
		...r,
		children: [t && (0, d.jsx)(`title`, { id: i, children: t }), pe(e)],
	});
}
function pe(e) {
	switch (e) {
		case `course`:
			return (0, d.jsxs)(d.Fragment, {
				children: [
					(0, d.jsx)(`rect`, { x: `2`, y: `1.5`, width: `11`, height: `12` }),
					(0, d.jsx)(`line`, { x1: `4.5`, y1: `5`, x2: `10.5`, y2: `5` }),
					(0, d.jsx)(`line`, { x1: `4.5`, y1: `8`, x2: `10.5`, y2: `8` }),
					(0, d.jsx)(`line`, { x1: `4.5`, y1: `11`, x2: `8`, y2: `11` }),
				],
			});
		case `slides`:
			return (0, d.jsxs)(d.Fragment, {
				children: [
					(0, d.jsx)(`rect`, { x: `1.5`, y: `2.5`, width: `12`, height: `8.5` }),
					(0, d.jsx)(`line`, { x1: `6`, y1: `13`, x2: `9`, y2: `13` }),
					(0, d.jsx)(`line`, { x1: `7.5`, y1: `11`, x2: `7.5`, y2: `13` }),
				],
			});
		case `visualize`:
			return (0, d.jsxs)(d.Fragment, {
				children: [
					(0, d.jsx)(`rect`, { className: de.filled, x: `1.5`, y: `4`, width: `3`, height: `7` }),
					(0, d.jsx)(`rect`, { x: `5.5`, y: `2`, width: `3`, height: `9` }),
					(0, d.jsx)(`rect`, { className: de.filled, x: `9.5`, y: `6`, width: `3`, height: `5` }),
				],
			});
		case `upload`:
			return (0, d.jsxs)(d.Fragment, {
				children: [
					(0, d.jsx)(`path`, { d: `M2 9 L2 13 L13 13 L13 9` }),
					(0, d.jsx)(`line`, { x1: `7.5`, y1: `2`, x2: `7.5`, y2: `9.5` }),
					(0, d.jsx)(`path`, { d: `M4.5 5 L7.5 2 L10.5 5` }),
				],
			});
		case `transcript`:
			return (0, d.jsxs)(d.Fragment, {
				children: [
					(0, d.jsx)(`rect`, { x: `1.5`, y: `2`, width: `12`, height: `11` }),
					(0, d.jsx)(`line`, { x1: `4`, y1: `5`, x2: `11`, y2: `5` }),
					(0, d.jsx)(`line`, { x1: `4`, y1: `7.5`, x2: `11`, y2: `7.5` }),
					(0, d.jsx)(`line`, { x1: `4`, y1: `10`, x2: `8`, y2: `10` }),
				],
			});
		case `comments`:
			return (0, d.jsx)(`path`, { d: `M2 2 L13 2 L13 9 L7 9 L4 12 L4 9 L2 9 Z` });
		case `handout`:
			return (0, d.jsxs)(d.Fragment, {
				children: [
					(0, d.jsx)(`path`, { d: `M3 1.5 L10 1.5 L12.5 4 L12.5 13.5 L3 13.5 Z` }),
					(0, d.jsx)(`path`, { d: `M10 1.5 L10 4 L12.5 4` }),
				],
			});
	}
}
var D = {
		root: `_root_n5d4d_1`,
		sizeBody: `_sizeBody_n5d4d_6`,
		sizeCompact: `_sizeCompact_n5d4d_10`,
		sizeMetadata: `_sizeMetadata_n5d4d_14`,
		sizeLabel: `_sizeLabel_n5d4d_18`,
		sizeMetric: `_sizeMetric_n5d4d_24`,
		tonePrimary: `_tonePrimary_n5d4d_28`,
		toneMuted: `_toneMuted_n5d4d_29`,
		toneInverse: `_toneInverse_n5d4d_30`,
		toneAccent: `_toneAccent_n5d4d_31`,
		toneSuccess: `_toneSuccess_n5d4d_32`,
		toneWarning: `_toneWarning_n5d4d_33`,
		toneDanger: `_toneDanger_n5d4d_34`,
		toneInherit: `_toneInherit_n5d4d_35`,
		weightRegular: `_weightRegular_n5d4d_37`,
		weightBold: `_weightBold_n5d4d_38`,
		alignStart: `_alignStart_n5d4d_40`,
		alignCenter: `_alignCenter_n5d4d_41`,
		alignEnd: `_alignEnd_n5d4d_42`,
		truncate: `_truncate_n5d4d_44`,
	},
	O = n(),
	k = {
		body: D.sizeBody,
		compact: D.sizeCompact,
		metadata: D.sizeMetadata,
		label: D.sizeLabel,
		metric: D.sizeMetric,
	},
	me = {
		primary: D.tonePrimary,
		muted: D.toneMuted,
		inverse: D.toneInverse,
		accent: D.toneAccent,
		success: D.toneSuccess,
		warning: D.toneWarning,
		danger: D.toneDanger,
		inherit: D.toneInherit,
	},
	he = { regular: D.weightRegular, bold: D.weightBold },
	ge = { start: D.alignStart, center: D.alignCenter, end: D.alignEnd };
function _e({
	as: e = `p`,
	size: t = `body`,
	tone: n = `primary`,
	weight: r = `regular`,
	align: i = `start`,
	truncate: a = !1,
	className: o,
	children: s,
	...c
}) {
	return (0, O.createElement)(
		e,
		{
			className: [D.root, k[t], me[n], he[r], ge[i], a ? D.truncate : ``, o ?? ``]
				.filter(Boolean)
				.join(` `),
			...c,
		},
		s,
	);
}
var ve = {
		root: `_root_rss1l_1`,
		inline: `_inline_rss1l_8`,
		block: `_block_rss1l_9`,
		muted: `_muted_rss1l_10`,
		accent: `_accent_rss1l_11`,
		copyable: `_copyable_rss1l_13`,
	},
	ye = { primary: ``, muted: ve.muted, accent: ve.accent };
function be({
	as: e = `code`,
	tone: t = `primary`,
	display: n = `inline`,
	copyable: r = !1,
	className: i,
	children: a,
	onClick: o,
	...s
}) {
	let c = [ve.root, n === `inline` ? ve.inline : ve.block, ye[t], r ? ve.copyable : ``, i ?? ``]
		.filter(Boolean)
		.join(` `);
	return (0, d.jsx)(e, {
		className: c,
		onClick: (e) => {
			r && typeof a == `string` && navigator.clipboard?.writeText(a).catch(() => {}), o?.(e);
		},
		...s,
		children: a,
	});
}
var xe = {
		root: `_root_2cwmi_1`,
		pending: `_pending_2cwmi_5`,
		ready: `_ready_2cwmi_6`,
		running: `_running_2cwmi_7`,
		succeeded: `_succeeded_2cwmi_8`,
		done: `_done_2cwmi_9`,
		partial: `_partial_2cwmi_10`,
		warning: `_warning_2cwmi_11`,
		failed: `_failed_2cwmi_12`,
		error: `_error_2cwmi_13`,
		canceled: `_canceled_2cwmi_14`,
	},
	Se = {
		pending: `◌`,
		ready: `◌`,
		running: `●`,
		succeeded: `✔`,
		done: `●`,
		partial: `◐`,
		warning: `⚠`,
		failed: `✘`,
		error: `✘`,
		canceled: `⊘`,
	};
function Ce({ status: e, icon: t = !1, className: n, children: r, ...i }) {
	let a = e.toLowerCase(),
		o = xe[a] ?? ``,
		s = r ?? e;
	return (0, d.jsxs)(`span`, {
		className: [xe.root, o, n ?? ``].filter(Boolean).join(` `),
		...i,
		children: [t ? `${Se[a] ?? `?`} ` : ``, s],
	});
}
var we = { root: `_root_1fiw6_1`, vertical: `_vertical_1fiw6_7` };
function Te({ orientation: e = `horizontal`, className: t, ...n }) {
	return (0, d.jsx)(`hr`, {
		"aria-orientation": e,
		className: [we.root, e === `vertical` ? we.vertical : ``, t ?? ``].filter(Boolean).join(` `),
		...n,
	});
}
var Ee = {
		root: `_root_kfqnk_1`,
		uppercase: `_uppercase_kfqnk_7`,
		truncate: `_truncate_kfqnk_8`,
		toneMuted: `_toneMuted_kfqnk_9`,
		toneAccent: `_toneAccent_kfqnk_10`,
		toneSuccess: `_toneSuccess_kfqnk_11`,
		toneWarning: `_toneWarning_kfqnk_12`,
		toneDanger: `_toneDanger_kfqnk_13`,
		toneInherit: `_toneInherit_kfqnk_14`,
	},
	De = {
		muted: Ee.toneMuted,
		accent: Ee.toneAccent,
		success: Ee.toneSuccess,
		warning: Ee.toneWarning,
		danger: Ee.toneDanger,
		inherit: Ee.toneInherit,
	};
function A({
	tone: e = `muted`,
	transform: t = `none`,
	truncate: n = !1,
	className: r,
	children: i,
	...a
}) {
	return (0, d.jsx)(`span`, {
		className: [
			Ee.root,
			De[e],
			t === `uppercase` ? Ee.uppercase : ``,
			n ? Ee.truncate : ``,
			r ?? ``,
		]
			.filter(Boolean)
			.join(` `),
		"data-rag-foundation": `Caption`,
		...a,
		children: i,
	});
}
var Oe = {
	root: `_root_ga2ee_1`,
	fill: `_fill_ga2ee_7`,
	header: `_header_ga2ee_11`,
	body: `_body_ga2ee_26`,
	condensed: `_condensed_ga2ee_30`,
};
function ke({
	title: e,
	actions: t,
	density: n = `normal`,
	fill: r = !1,
	className: i,
	children: a,
	...o
}) {
	return (0, d.jsxs)(`section`, {
		className: [Oe.root, r ? Oe.fill : ``, i ?? ``].filter(Boolean).join(` `),
		"data-rag-layout": `Panel`,
		...o,
		children: [
			(e || t) &&
				(0, d.jsxs)(`div`, {
					className: Oe.header,
					children: [(0, d.jsx)(`span`, { children: e }), t],
				}),
			(0, d.jsx)(`div`, { className: n === `condensed` ? Oe.condensed : Oe.body, children: a }),
		],
	});
}
var Ae = {
		root: `_root_15awp_1`,
		gapXs: `_gapXs_15awp_6`,
		gapSm: `_gapSm_15awp_7`,
		gapMd: `_gapMd_15awp_8`,
		gapLg: `_gapLg_15awp_9`,
		gapXl: `_gapXl_15awp_10`,
		alignStretch: `_alignStretch_15awp_12`,
		alignStart: `_alignStart_15awp_13`,
		alignCenter: `_alignCenter_15awp_14`,
		alignEnd: `_alignEnd_15awp_15`,
	},
	je = { xs: Ae.gapXs, sm: Ae.gapSm, md: Ae.gapMd, lg: Ae.gapLg, xl: Ae.gapXl },
	Me = { stretch: Ae.alignStretch, start: Ae.alignStart, center: Ae.alignCenter, end: Ae.alignEnd };
function Ne({ gap: e = `md`, align: t = `stretch`, className: n, children: r, ...i }) {
	return (0, d.jsx)(`div`, {
		className: [Ae.root, je[e], Me[t], n ?? ``].filter(Boolean).join(` `),
		"data-rag-layout": `Stack`,
		...i,
		children: r,
	});
}
var Pe = {
		root: `_root_8hatg_1`,
		wrap: `_wrap_8hatg_6`,
		noWrap: `_noWrap_8hatg_7`,
		gapXs: `_gapXs_8hatg_8`,
		gapSm: `_gapSm_8hatg_9`,
		gapMd: `_gapMd_8hatg_10`,
		gapLg: `_gapLg_8hatg_11`,
		justifyStart: `_justifyStart_8hatg_12`,
		justifyBetween: `_justifyBetween_8hatg_13`,
		justifyEnd: `_justifyEnd_8hatg_14`,
	},
	Fe = { xs: Pe.gapXs, sm: Pe.gapSm, md: Pe.gapMd, lg: Pe.gapLg },
	Ie = { start: Pe.justifyStart, between: Pe.justifyBetween, end: Pe.justifyEnd };
function Le({
	gap: e = `sm`,
	justify: t = `start`,
	wrap: n = !0,
	className: r,
	children: i,
	...a
}) {
	return (0, d.jsx)(`div`, {
		className: [Pe.root, Fe[e], Ie[t], n ? Pe.wrap : Pe.noWrap, r ?? ``].filter(Boolean).join(` `),
		"data-rag-layout": `Inline`,
		...a,
		children: i,
	});
}
var Re = {
		root: `_root_rmtbr_1`,
		searchWorkbench: `_searchWorkbench_rmtbr_7`,
		corpusExplorer: `_corpusExplorer_rmtbr_12`,
		twoColumn: `_twoColumn_rmtbr_17`,
	},
	ze = {
		searchWorkbench: Re.searchWorkbench,
		corpusExplorer: Re.corpusExplorer,
		twoColumn: Re.twoColumn,
	};
function Be({ recipe: e = `twoColumn`, className: t, children: n, ...r }) {
	return (0, d.jsx)(`div`, {
		className: [Re.root, ze[e], t ?? ``].filter(Boolean).join(` `),
		"data-rag-layout": `DashboardGrid`,
		"data-rag-layout-recipe": e,
		...r,
		children: n,
	});
}
var Ve = {
	root: `_root_1rq8r_1`,
	header: `_header_1rq8r_2`,
	body: `_body_1rq8r_3`,
	sidebar: `_sidebar_1rq8r_4`,
	main: `_main_1rq8r_5`,
};
function He({ header: e, sidebar: t, className: n, children: r, ...i }) {
	return (0, d.jsxs)(`div`, {
		className: [Ve.root, n ?? ``].filter(Boolean).join(` `),
		"data-rag-layout": `AppShell`,
		...i,
		children: [
			e && (0, d.jsx)(`header`, { className: Ve.header, children: e }),
			(0, d.jsxs)(`div`, {
				className: Ve.body,
				children: [
					t && (0, d.jsx)(`aside`, { className: Ve.sidebar, children: t }),
					(0, d.jsx)(`main`, { className: Ve.main, children: r }),
				],
			}),
		],
	});
}
var Ue = { root: `_root_37934_1`, y: `_y_37934_2`, x: `_x_37934_3`, both: `_both_37934_4` };
function We({ axis: e = `y`, className: t, children: n, ...r }) {
	return (0, d.jsx)(`div`, {
		className: [Ue.root, Ue[e], t ?? ``].filter(Boolean).join(` `),
		"data-rag-layout": `ScrollRegion`,
		...r,
		children: n,
	});
}
var Ge = { root: `_root_bo8as_1`, item: `_item_bo8as_2`, active: `_active_bo8as_4` };
function Ke({ items: e, activeId: t, onChange: n, ariaLabel: r = `Tabs` }) {
	return (0, d.jsx)(`div`, {
		className: Ge.root,
		role: `tablist`,
		"aria-label": r,
		"data-rag-layout": `TabList`,
		children: e.map((e) =>
			(0, d.jsx)(
				`button`,
				{
					type: `button`,
					role: `tab`,
					"aria-selected": e.id === t,
					className: [Ge.item, e.id === t ? Ge.active : ``].filter(Boolean).join(` `),
					onClick: () => n(e.id),
					children: e.label,
				},
				e.id,
			),
		),
	});
}
var qe = {
	root: `_root_t0l9w_1`,
	label: `_label_t0l9w_2`,
	control: `_control_t0l9w_3`,
	hint: `_hint_t0l9w_4`,
};
function Je({ label: e, control: t, hint: n, className: r, ...i }) {
	return (0, d.jsxs)(`div`, {
		className: [qe.root, r ?? ``].filter(Boolean).join(` `),
		"data-rag-layout": `FormRow`,
		...i,
		children: [
			(0, d.jsx)(`div`, { className: qe.label, children: e }),
			(0, d.jsx)(`div`, { className: qe.control, children: t }),
			n && (0, d.jsx)(`div`, { className: qe.hint, children: n }),
		],
	});
}
var Ye = {
	root: `_root_145df_1`,
	normal: `_normal_145df_2`,
	spacious: `_spacious_145df_3`,
	"divider-top": `_divider-top_145df_4`,
	"divider-bottom": `_divider-bottom_145df_5`,
	"divider-both": `_divider-both_145df_6`,
	label: `_label_145df_7`,
	caption: `_caption_145df_8`,
	body: `_body_145df_9`,
};
function Xe({
	as: e = `section`,
	label: t,
	caption: n,
	density: r = `normal`,
	divider: i = `none`,
	className: a,
	children: o,
	...s
}) {
	return (0, d.jsxs)(e, {
		className: [Ye.root, Ye[r], Ye[`divider-${i}`], a ?? ``].filter(Boolean).join(` `),
		"data-rag-layout": `SectionBlock`,
		...s,
		children: [
			t && (0, d.jsx)(`div`, { className: Ye.label, children: t }),
			n && (0, d.jsx)(A, { className: Ye.caption, children: n }),
			(0, d.jsx)(`div`, { className: Ye.body, children: o }),
		],
	});
}
var Ze = {
	root: `_root_1o5gv_1`,
	balanced: `_balanced_1o5gv_2`,
	leftNarrow: `_leftNarrow_1o5gv_3`,
	rightNarrow: `_rightNarrow_1o5gv_4`,
	course: `_course_1o5gv_5`,
	sidebar: `_sidebar_1o5gv_6`,
	divider: `_divider_1o5gv_7`,
	pane: `_pane_1o5gv_7`,
	gutterMd: `_gutterMd_1o5gv_9`,
	gutterLg: `_gutterLg_1o5gv_11`,
};
function Qe({
	left: e,
	right: t,
	ratio: n = `balanced`,
	divider: r = !1,
	gutter: i = `none`,
	className: a,
	...o
}) {
	return (0, d.jsxs)(`div`, {
		className: [
			Ze.root,
			Ze[n],
			r ? Ze.divider : ``,
			i === `none` ? `` : Ze[`gutter${i[0].toUpperCase()}${i.slice(1)}`],
			a ?? ``,
		]
			.filter(Boolean)
			.join(` `),
		"data-rag-layout": `SplitPane`,
		...o,
		children: [
			(0, d.jsx)(`div`, { className: Ze.pane, children: e }),
			(0, d.jsx)(`div`, { className: Ze.pane, children: t }),
		],
	});
}
var $e = {
	root: `_root_106h4_1`,
	header: `_header_106h4_12`,
	metaRow: `_metaRow_106h4_20`,
	title: `_title_106h4_27`,
	subtitle: `_subtitle_106h4_35`,
	body: `_body_106h4_42`,
	twoPane: `_twoPane_106h4_47`,
	singlePane: `_singlePane_106h4_52`,
	primaryWide: `_primaryWide_106h4_56`,
	"primary-left": `_primary-left_106h4_56`,
	"primary-right": `_primary-right_106h4_57`,
	secondaryWide: `_secondaryWide_106h4_61`,
	balanced: `_balanced_106h4_66`,
	primary: `_primary_106h4_56`,
	secondary: `_secondary_106h4_61`,
	divider: `_divider_106h4_90`,
	footer: `_footer_106h4_96`,
};
function et({
	as: e = `article`,
	eyebrow: t,
	counter: n,
	title: r,
	subtitle: i,
	primary: a,
	secondary: o,
	primarySide: s = `left`,
	ratio: c = `primaryWide`,
	divider: l = !0,
	footer: u,
	className: f,
	...p
}) {
	return (0, d.jsxs)(e, {
		className: [$e.root, f ?? ``].filter(Boolean).join(` `),
		"data-rag-layout": `SlideShell`,
		"data-primary-side": s,
		...p,
		children: [
			(0, d.jsxs)(`header`, {
				className: $e.header,
				children: [
					(0, d.jsxs)(`div`, {
						className: $e.metaRow,
						children: [
							t ? (0, d.jsx)(A, { transform: `uppercase`, children: t }) : (0, d.jsx)(`span`, {}),
							n && (0, d.jsx)(A, { children: n }),
						],
					}),
					(0, d.jsx)(`h2`, { className: $e.title, children: r }),
					i && (0, d.jsx)(`p`, { className: $e.subtitle, children: i }),
				],
			}),
			(0, d.jsxs)(`div`, {
				className: [
					$e.body,
					$e[c],
					o ? $e.twoPane : $e.singlePane,
					l && o ? $e.divider : ``,
					$e[`primary-${s}`],
				]
					.filter(Boolean)
					.join(` `),
				children: [
					(0, d.jsx)(`div`, { className: $e.primary, children: a }),
					o && (0, d.jsx)(`div`, { className: $e.secondary, children: o }),
				],
			}),
			u && (0, d.jsx)(`footer`, { className: $e.footer, children: u }),
		],
	});
}
var tt = {
	root: `_root_1w0ju_1`,
	sidebar: `_sidebar_1w0ju_10`,
	sidebarHeader: `_sidebarHeader_1w0ju_19`,
	sidebarFooter: `_sidebarFooter_1w0ju_20`,
	sidebarBody: `_sidebarBody_1w0ju_33`,
	content: `_content_1w0ju_40`,
	contentPaddingMd: `_contentPaddingMd_1w0ju_47`,
	contentPaddingLg: `_contentPaddingLg_1w0ju_51`,
};
function nt(e) {
	return typeof e == `number` ? `${e}px` : e;
}
function rt({
	sidebar: e,
	sidebarWidth: t = 188,
	contentPadding: n = `none`,
	header: r,
	footer: i,
	className: a,
	children: o,
	style: s,
	...c
}) {
	let l = { ...s, "--rag-sidebar-width": nt(t) };
	return (0, d.jsxs)(`div`, {
		className: [tt.root, a ?? ``].filter(Boolean).join(` `),
		style: l,
		"data-rag-layout": `SidebarShell`,
		...c,
		children: [
			(0, d.jsxs)(`aside`, {
				className: tt.sidebar,
				"aria-label": `Course navigation`,
				children: [
					r && (0, d.jsx)(`div`, { className: tt.sidebarHeader, children: r }),
					(0, d.jsx)(`div`, { className: tt.sidebarBody, children: e }),
					i && (0, d.jsx)(`div`, { className: tt.sidebarFooter, children: i }),
				],
			}),
			(0, d.jsx)(`main`, {
				className: [
					tt.content,
					n === `none` ? `` : tt[`contentPadding${n[0].toUpperCase()}${n.slice(1)}`],
				]
					.filter(Boolean)
					.join(` `),
				children: o,
			}),
		],
	});
}
var it = {
	root: `_root_1gyp0_1`,
	row: `_row_1gyp0_2`,
	key: `_key_1gyp0_3`,
	value: `_value_1gyp0_4`,
	compact: `_compact_1gyp0_5`,
	copyButton: `_copyButton_1gyp0_6`,
};
function at({ items: e, density: t = `normal`, className: n, onCopy: r }) {
	let i = (e) => {
		if (r) {
			r(e);
			return;
		}
		navigator.clipboard?.writeText(e).catch(() => {});
	};
	return (0, d.jsx)(`dl`, {
		className: [it.root, t === `compact` ? it.compact : ``, n ?? ``].filter(Boolean).join(` `),
		"data-rag-component": `MetadataGrid`,
		children: e.map((e, t) =>
			(0, d.jsxs)(
				`div`,
				{
					className: it.row,
					children: [
						(0, d.jsx)(`dt`, { className: it.key, children: e.key }),
						(0, d.jsxs)(`dd`, {
							className: it.value,
							children: [
								e.value,
								e.copyValue &&
									(0, d.jsx)(`button`, {
										type: `button`,
										className: it.copyButton,
										onClick: () => i(e.copyValue),
										children: `⧉`,
									}),
							],
						}),
					],
				},
				t,
			),
		),
	});
}
var ot = {
	root: `_root_1kw0z_1`,
	selectable: `_selectable_1kw0z_4`,
	selected: `_selected_1kw0z_6`,
	start: `_start_1kw0z_7`,
	end: `_end_1kw0z_8`,
	center: `_center_1kw0z_9`,
	empty: `_empty_1kw0z_10`,
};
function st({
	columns: e,
	rows: t,
	getRowKey: n,
	selectedKey: r,
	onRowSelect: i,
	emptyMessage: a,
	className: o,
}) {
	return (0, d.jsxs)(`table`, {
		className: [ot.root, o ?? ``].filter(Boolean).join(` `),
		"data-rag-component": `DataTable`,
		children: [
			(0, d.jsx)(`thead`, {
				children: (0, d.jsx)(`tr`, {
					children: e.map((e) =>
						(0, d.jsx)(`th`, { className: ot[e.align ?? `start`], children: e.header }, e.id),
					),
				}),
			}),
			(0, d.jsxs)(`tbody`, {
				children: [
					t.length === 0 &&
						a &&
						(0, d.jsx)(`tr`, {
							children: (0, d.jsx)(`td`, { colSpan: e.length, className: ot.empty, children: a }),
						}),
					t.map((t) => {
						let a = n(t);
						return (0, d.jsx)(
							`tr`,
							{
								className: [i ? ot.selectable : ``, r === a ? ot.selected : ``]
									.filter(Boolean)
									.join(` `),
								onClick: i ? () => i(t) : void 0,
								children: e.map((e) =>
									(0, d.jsx)(
										`td`,
										{
											className: ot[e.align ?? `start`],
											style: e.maxWidth ? { maxWidth: e.maxWidth } : void 0,
											children: e.cell(t),
										},
										e.id,
									),
								),
							},
							a,
						);
					}),
				],
			}),
		],
	});
}
var ct = {
	root: `_root_13gif_1`,
	brand: `_brand_13gif_13`,
	item: `_item_13gif_22`,
	active: `_active_13gif_42`,
};
function lt({ brand: e, items: t, activeItemId: n, onItemSelect: r }) {
	return (0, d.jsxs)(`nav`, {
		className: ct.root,
		"aria-label": `Primary`,
		"data-rag-component": `AppNav`,
		children: [
			(0, d.jsx)(`span`, { className: ct.brand, children: e }),
			t.map((e) => {
				let t = e.id === n;
				return (0, d.jsx)(
					`button`,
					{
						type: `button`,
						className: [ct.item, t ? ct.active : ``].filter(Boolean).join(` `),
						"aria-current": t ? `page` : void 0,
						onClick: () => r(e.id),
						children: e.label,
					},
					e.id,
				);
			}),
		],
	});
}
var ut = {
	root: `_root_1i26l_1`,
	default: `_default_1i26l_7`,
	compact: `_compact_1i26l_11`,
	item: `_item_1i26l_15`,
	label: `_label_1i26l_22`,
};
function dt({
	kinds: e = x.filter((e) => e !== `other`),
	mode: t = `pattern`,
	compact: n = !1,
	selectedKind: r,
	className: i,
	...a
}) {
	return (0, d.jsx)(`div`, {
		className: [ut.root, n ? ut.compact : ut.default, i ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `ContextLegend`,
		"data-mode": t,
		...a,
		children: e.map((e) =>
			(0, d.jsxs)(
				`div`,
				{
					className: ut.item,
					"data-selected": r === e ? `true` : void 0,
					children: [
						(0, d.jsx)(oe, { kind: e, mode: t, size: n ? `sm` : `md`, selected: r === e }),
						(0, d.jsx)(`span`, { className: ut.label, children: ne(e) }),
					],
				},
				e,
			),
		),
	});
}
var ft = {
	root: `_root_4isez_1`,
	header: `_header_4isez_8`,
	track: `_track_4isez_15`,
	segment: `_segment_4isez_25`,
	selected: `_selected_4isez_32`,
	limitMarker: `_limitMarker_4isez_38`,
	nearBudget: `_nearBudget_4isez_47`,
	overBudget: `_overBudget_4isez_51`,
	legend: `_legend_4isez_55`,
	legendItem: `_legendItem_4isez_61`,
	dot: `_dot_4isez_74`,
	tokens: `_tokens_4isez_81`,
	kind_system: `_kind_system_4isez_85`,
	kind_instruction: `_kind_instruction_4isez_86`,
	kind_context: `_kind_context_4isez_87`,
	kind_conversation: `_kind_conversation_4isez_88`,
	kind_course: `_kind_course_4isez_89`,
	kind_other: `_kind_other_4isez_90`,
	kind_summary: `_kind_summary_4isez_91`,
	kind_annotation: `_kind_annotation_4isez_92`,
	kind_retrieval: `_kind_retrieval_4isez_93`,
	kind_tool: `_kind_tool_4isez_94`,
	kind_result: `_kind_result_4isez_95`,
	kind_generated: `_kind_generated_4isez_96`,
	kind_evicted: `_kind_evicted_4isez_97`,
	kind_active: `_kind_active_4isez_98`,
	kind_empty: `_kind_empty_4isez_99`,
};
function pt(e) {
	return `${Math.round(e).toLocaleString()} tok`;
}
function mt(e) {
	return e.filter((e) => e.kind !== `empty` && e.tokens > 0);
}
function ht({
	snapshot: e,
	mode: t = `pattern`,
	showLegend: n = !0,
	selectedPartId: r,
	className: i,
	...a
}) {
	let o = w(e),
		s = ie(e),
		c = o > e.limit,
		l = !c && s >= 0.8,
		u = mt(e.parts),
		f = r ?? e.selectedPartId;
	return (0, d.jsxs)(`div`, {
		className: [ft.root, c ? ft.overBudget : l ? ft.nearBudget : ``, i ?? ``]
			.filter(Boolean)
			.join(` `),
		"data-rag-molecule": `ContextBudgetBar`,
		"data-mode": t,
		...a,
		children: [
			(0, d.jsxs)(`div`, {
				className: ft.header,
				children: [
					(0, d.jsx)(A, { transform: `uppercase`, children: e.title }),
					(0, d.jsxs)(A, {
						tone: c ? `danger` : l ? `warning` : `muted`,
						children: [pt(o), ` / `, pt(e.limit), ` · `, Math.round(s * 100), `%`],
					}),
				],
			}),
			(0, d.jsxs)(`div`, {
				className: ft.track,
				role: `img`,
				"aria-label": `${e.title}: ${pt(o)} of ${pt(e.limit)} used`,
				children: [
					u.map((t) => {
						let n = e.limit > 0 ? Math.max(0.5, (t.tokens / e.limit) * 100) : 0;
						return (0, d.jsx)(
							`div`,
							{
								className: [
									ft.segment,
									ft[`kind_${t.kind}`] ?? ft.kind_other,
									f === t.id ? ft.selected : ``,
								]
									.filter(Boolean)
									.join(` `),
								style: { width: `${n}%` },
								title: `${t.label}: ${pt(t.tokens)} (${ne(t.kind)})`,
							},
							t.id,
						);
					}),
					c && (0, d.jsx)(`div`, { className: ft.limitMarker, title: `model context limit` }),
				],
			}),
			n &&
				(0, d.jsx)(`div`, {
					className: ft.legend,
					children: u.map((e) =>
						(0, d.jsxs)(
							`span`,
							{
								className: ft.legendItem,
								"data-selected": f === e.id ? `true` : void 0,
								children: [
									(0, d.jsx)(`span`, {
										className: [ft.dot, ft[`kind_${e.kind}`] ?? ft.kind_other].join(` `),
									}),
									(0, d.jsx)(`span`, { children: e.label }),
									(0, d.jsx)(`span`, { className: ft.tokens, children: pt(e.tokens) }),
								],
							},
							e.id,
						),
					),
				}),
		],
	});
}
var gt = {
	root: `_root_hloz0_1`,
	strip: `_strip_hloz0_2`,
	segment: `_segment_hloz0_3`,
	label: `_label_hloz0_5`,
	selected: `_selected_hloz0_6`,
	caption: `_caption_hloz0_7`,
	kind_system: `_kind_system_hloz0_8`,
	kind_instruction: `_kind_instruction_hloz0_8`,
	kind_context: `_kind_context_hloz0_9`,
	kind_conversation: `_kind_conversation_hloz0_9`,
	kind_course: `_kind_course_hloz0_9`,
	kind_other: `_kind_other_hloz0_9`,
	kind_summary: `_kind_summary_hloz0_10`,
	kind_annotation: `_kind_annotation_hloz0_10`,
	kind_retrieval: `_kind_retrieval_hloz0_11`,
	kind_tool: `_kind_tool_hloz0_11`,
	kind_result: `_kind_result_hloz0_11`,
	kind_generated: `_kind_generated_hloz0_12`,
	kind_evicted: `_kind_evicted_hloz0_13`,
	kind_active: `_kind_active_hloz0_15`,
	kind_empty: `_kind_empty_hloz0_17`,
};
function _t(e) {
	return `${e.toLocaleString()} tok`;
}
function vt({
	snapshot: e,
	mode: t = `pattern`,
	selectedPartId: n,
	showLabels: r = !0,
	className: i,
	...a
}) {
	let o = n ?? e.selectedPartId,
		s = e.limit || e.parts.reduce((e, t) => e + t.tokens, 0);
	return (0, d.jsxs)(`div`, {
		className: [gt.root, i ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `ContextStripDiagram`,
		"data-mode": t,
		...a,
		children: [
			(0, d.jsx)(`div`, {
				className: gt.strip,
				role: `img`,
				"aria-label": `${e.title} context strip`,
				children: e.parts.map((e) => {
					let t = Math.max(2, (e.tokens / s) * 100),
						n = o === e.id;
					return (0, d.jsx)(
						`div`,
						{
							className: [gt.segment, gt[`kind_${e.kind}`] ?? gt.kind_other, n ? gt.selected : ``]
								.filter(Boolean)
								.join(` `),
							style: { width: `${t}%` },
							title: `${e.label}: ${_t(e.tokens)} (${ne(e.kind)})`,
							children:
								r && t >= 7 && (0, d.jsx)(`span`, { className: gt.label, children: e.label }),
						},
						e.id,
					);
				}),
			}),
			(0, d.jsxs)(`div`, {
				className: gt.caption,
				children: [e.title, ` · `, e.limit.toLocaleString(), ` token window`],
			}),
		],
	});
}
var yt = {
	root: `_root_eoqke_1`,
	layers: `_layers_eoqke_2`,
	layer: `_layer_eoqke_2`,
	label: `_label_eoqke_5`,
	tokens: `_tokens_eoqke_6`,
	kind_system: `_kind_system_eoqke_8`,
	kind_instruction: `_kind_instruction_eoqke_10`,
	kind_summary: `_kind_summary_eoqke_12`,
	kind_annotation: `_kind_annotation_eoqke_14`,
	kind_retrieval: `_kind_retrieval_eoqke_16`,
	kind_tool: `_kind_tool_eoqke_18`,
	kind_result: `_kind_result_eoqke_20`,
	kind_generated: `_kind_generated_eoqke_22`,
	selected: `_selected_eoqke_29`,
	kind_context: `_kind_context_eoqke_31`,
	kind_conversation: `_kind_conversation_eoqke_31`,
	kind_course: `_kind_course_eoqke_31`,
	kind_other: `_kind_other_eoqke_31`,
	kind_evicted: `_kind_evicted_eoqke_35`,
	kind_active: `_kind_active_eoqke_37`,
	kind_empty: `_kind_empty_eoqke_39`,
	caption: `_caption_eoqke_40`,
};
function bt(e) {
	return `${e.toLocaleString()} tok`;
}
function xt({ snapshot: e, selectedPartId: t, className: n, ...r }) {
	let i = t ?? e.selectedPartId,
		a = e.parts.filter((e) => e.tokens > 0),
		o = Math.max(...a.map((e) => e.tokens), 1);
	return (0, d.jsxs)(`div`, {
		className: [yt.root, n ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `ContextStackDiagram`,
		...r,
		children: [
			(0, d.jsx)(`div`, {
				className: yt.layers,
				role: `img`,
				"aria-label": `${e.title} context stack`,
				children: a.map((e) => {
					let t = i === e.id,
						n = Math.max(26, Math.min(74, 22 + (e.tokens / o) * 52));
					return (0, d.jsxs)(
						`div`,
						{
							className: [yt.layer, yt[`kind_${e.kind}`] ?? yt.kind_other, t ? yt.selected : ``]
								.filter(Boolean)
								.join(` `),
							style: { minHeight: n },
							title: `${e.label}: ${bt(e.tokens)} (${ne(e.kind)})`,
							children: [
								(0, d.jsx)(`span`, { className: yt.label, children: e.label }),
								(0, d.jsx)(`span`, { className: yt.tokens, children: bt(e.tokens) }),
							],
						},
						e.id,
					);
				}),
			}),
			(0, d.jsx)(`div`, { className: yt.caption, children: e.title }),
		],
	});
}
var St = {
	root: `_root_shrue_1`,
	map: `_map_shrue_2`,
	tile: `_tile_shrue_3`,
	label: `_label_shrue_4`,
	tokens: `_tokens_shrue_5`,
	kind_system: `_kind_system_shrue_7`,
	kind_instruction: `_kind_instruction_shrue_9`,
	kind_summary: `_kind_summary_shrue_11`,
	kind_annotation: `_kind_annotation_shrue_13`,
	kind_retrieval: `_kind_retrieval_shrue_15`,
	kind_tool: `_kind_tool_shrue_17`,
	kind_result: `_kind_result_shrue_19`,
	kind_generated: `_kind_generated_shrue_21`,
	selected: `_selected_shrue_29`,
	kind_context: `_kind_context_shrue_31`,
	kind_conversation: `_kind_conversation_shrue_31`,
	kind_course: `_kind_course_shrue_31`,
	kind_other: `_kind_other_shrue_31`,
	kind_evicted: `_kind_evicted_shrue_35`,
	kind_active: `_kind_active_shrue_37`,
	caption: `_caption_shrue_39`,
};
function Ct(e) {
	return `${e.toLocaleString()} tok`;
}
function wt({ snapshot: e, selectedPartId: t, className: n, ...r }) {
	let i = t ?? e.selectedPartId,
		a = e.parts.filter((e) => e.kind !== `empty` && e.tokens > 0),
		o = w(e) || 1;
	return (0, d.jsxs)(`div`, {
		className: [St.root, n ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `ContextTreemap`,
		...r,
		children: [
			(0, d.jsx)(`div`, {
				className: St.map,
				role: `img`,
				"aria-label": `${e.title} token treemap`,
				children: a.map((e) => {
					let t = Math.max(7, (e.tokens / o) * 100),
						n = i === e.id;
					return (0, d.jsxs)(
						`div`,
						{
							className: [St.tile, St[`kind_${e.kind}`] ?? St.kind_other, n ? St.selected : ``]
								.filter(Boolean)
								.join(` `),
							style: { flexBasis: `${t}%`, flexGrow: e.tokens },
							title: `${e.label}: ${Ct(e.tokens)} (${ne(e.kind)})`,
							children: [
								(0, d.jsx)(`span`, { className: St.label, children: e.label }),
								(0, d.jsx)(`span`, { className: St.tokens, children: Ct(e.tokens) }),
							],
						},
						e.id,
					);
				}),
			}),
			(0, d.jsxs)(`div`, { className: St.caption, children: [Ct(o), ` in use · `, e.title] }),
		],
	});
}
var Tt = {
		root: `_root_15h6i_1`,
		selected: `_selected_15h6i_10`,
		titleBar: `_titleBar_15h6i_15`,
		titleCluster: `_titleCluster_15h6i_26`,
		metaCluster: `_metaCluster_15h6i_27`,
		roleGlyph: `_roleGlyph_15h6i_38`,
		roleLabel: `_roleLabel_15h6i_43`,
		body: `_body_15h6i_51`,
		toolBody: `_toolBody_15h6i_58`,
		noteChip: `_noteChip_15h6i_64`,
		user: `_user_15h6i_80`,
		assistant: `_assistant_15h6i_90`,
		tool: `_tool_15h6i_58`,
		system: `_system_15h6i_92`,
		developer: `_developer_15h6i_93`,
		other: `_other_15h6i_94`,
	},
	Et = {
		system: `SYSTEM`,
		developer: `DEVELOPER`,
		user: `USER`,
		assistant: `ASSISTANT`,
		tool: `TOOL`,
		other: `MESSAGE`,
	},
	Dt = { system: `◆`, developer: `◇`, user: `▸`, assistant: `✦`, tool: `⚙`, other: `•` };
function Ot(e, t) {
	return e.findIndex((e) => e.id === t) + 1;
}
function kt({
	message: e,
	annotations: t = [],
	selectedAnnotationId: n,
	onAnnotationSelect: r,
	showAnnotationChips: i = !0,
	className: a,
	...o
}) {
	let s = t.filter((t) => t.targetMessageId === e.id || e.annotationIds?.includes(t.id)),
		c = s.some((e) => e.id === n),
		l = e.name ? `${Et[e.role]} · ${e.name}` : Et[e.role],
		u = [Tt.body, e.role === `tool` ? Tt.toolBody : ``].filter(Boolean).join(` `);
	return (0, d.jsxs)(`article`, {
		className: [Tt.root, Tt[e.role] ?? Tt.other, c ? Tt.selected : ``, a ?? ``]
			.filter(Boolean)
			.join(` `),
		"data-rag-molecule": `TranscriptMessageCard`,
		"data-rag-message-id": e.id,
		"data-rag-role": e.role,
		"data-rag-has-notes": s.length > 0 ? `true` : void 0,
		...o,
		children: [
			(0, d.jsxs)(`header`, {
				className: Tt.titleBar,
				children: [
					(0, d.jsxs)(`div`, {
						className: Tt.titleCluster,
						children: [
							(0, d.jsx)(`span`, {
								className: Tt.roleGlyph,
								"aria-hidden": `true`,
								children: Dt[e.role],
							}),
							(0, d.jsx)(`span`, { className: Tt.roleLabel, children: l }),
						],
					}),
					(0, d.jsxs)(`div`, {
						className: Tt.metaCluster,
						children: [
							i &&
								s.map((e) => {
									let i = Ot(t, e.id),
										a = n === e.id;
									return (0, d.jsxs)(
										`button`,
										{
											type: `button`,
											className: Tt.noteChip,
											"data-selected": a ? `true` : void 0,
											onClick: () => r?.(e.id),
											children: [`note `, i > 0 ? i : ``],
										},
										e.id,
									);
								}),
							(0, d.jsxs)(A, { children: [e.tokens?.toLocaleString() ?? `—`, ` tok`] }),
						],
					}),
				],
			}),
			(0, d.jsx)(`div`, { className: u, children: e.text }),
		],
	});
}
var At = {
	root: `_root_1q2ab_1`,
	titleRow: `_titleRow_1q2ab_11`,
	metaRow: `_metaRow_1q2ab_12`,
	rightSlot: `_rightSlot_1q2ab_20`,
};
function jt({
	title: e,
	subtitle: t,
	messageCount: n,
	annotationCount: r,
	tokenTotal: i,
	rightSlot: a,
	className: o,
	...s
}) {
	let c = [
		n == null ? void 0 : `${n} message${n === 1 ? `` : `s`}`,
		r == null ? void 0 : `${r} annotation${r === 1 ? `` : `s`}`,
	]
		.filter(Boolean)
		.join(` · `);
	return (0, d.jsxs)(`header`, {
		className: [At.root, o ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `TranscriptSessionHeader`,
		...s,
		children: [
			(0, d.jsxs)(`div`, {
				className: At.titleRow,
				children: [
					(0, d.jsxs)(_e, {
						as: `div`,
						size: `metric`,
						weight: `bold`,
						children: [`Session: `, e],
					}),
					i != null && (0, d.jsxs)(A, { children: [i.toLocaleString(), ` tok cumulative`] }),
				],
			}),
			(t || c || a) &&
				(0, d.jsxs)(`div`, {
					className: At.metaRow,
					children: [
						(0, d.jsx)(A, { children: [t, c].filter(Boolean).join(` · `) }),
						a && (0, d.jsx)(`div`, { className: At.rightSlot, children: a }),
					],
				}),
		],
	});
}
var Mt = {
	root: `_root_1hxft_1`,
	selected: `_selected_1hxft_9`,
	titleBar: `_titleBar_1hxft_14`,
	noteLabel: `_noteLabel_1hxft_29`,
	body: `_body_1hxft_34`,
	title: `_title_1hxft_14`,
	text: `_text_1hxft_45`,
};
function Nt({ annotation: e, selected: t = !1, index: n, className: r, ...i }) {
	return (0, d.jsxs)(`aside`, {
		className: [Mt.root, t ? Mt.selected : ``, r ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `AnnotationNoteCard`,
		"data-rag-note-id": e.id,
		"data-rag-target-message-id": e.targetMessageId,
		...i,
		children: [
			(0, d.jsxs)(`header`, {
				className: Mt.titleBar,
				children: [
					(0, d.jsxs)(`span`, { className: Mt.noteLabel, children: [`NOTE `, n ?? ``] }),
					(0, d.jsx)(oe, { kind: e.kind, size: `sm`, selected: t }),
				],
			}),
			(0, d.jsxs)(`div`, {
				className: Mt.body,
				children: [
					(0, d.jsx)(`div`, { className: Mt.title, children: e.label }),
					(0, d.jsx)(`p`, { className: Mt.text, children: e.text }),
					e.confidence != null &&
						(0, d.jsxs)(A, { children: [`confidence `, Math.round(e.confidence * 100), `%`] }),
				],
			}),
		],
	});
}
var Pt = {
	root: `_root_1c1xz_1`,
	selected: `_selected_1c1xz_2`,
	resolved: `_resolved_1c1xz_3`,
	header: `_header_1c1xz_4`,
	index: `_index_1c1xz_6`,
	author: `_author_1c1xz_7`,
	spacer: `_spacer_1c1xz_8`,
	time: `_time_1c1xz_9`,
	dismiss: `_dismiss_1c1xz_10`,
	text: `_text_1c1xz_11`,
	compact: `_compact_1c1xz_12`,
};
function Ft({
	comment: e,
	index: t,
	selected: n = !1,
	compact: r = !1,
	onDismiss: i,
	className: a,
	...o
}) {
	return (0, d.jsxs)(`article`, {
		className: [
			Pt.root,
			n ? Pt.selected : ``,
			r ? Pt.compact : ``,
			e.status === `resolved` ? Pt.resolved : ``,
			a ?? ``,
		]
			.filter(Boolean)
			.join(` `),
		"data-rag-molecule": `AnchoredCommentCard`,
		"data-status": e.status ?? `open`,
		...o,
		children: [
			(0, d.jsxs)(`header`, {
				className: Pt.header,
				children: [
					t != null && (0, d.jsx)(`span`, { className: Pt.index, children: t }),
					(0, d.jsx)(`span`, { className: Pt.author, children: e.author.toUpperCase() }),
					(0, d.jsx)(`span`, { className: Pt.spacer }),
					e.time && (0, d.jsx)(`span`, { className: Pt.time, children: e.time }),
					i &&
						(0, d.jsx)(`button`, {
							type: `button`,
							className: Pt.dismiss,
							onClick: i,
							"aria-label": `Dismiss comment`,
							children: `✕`,
						}),
				],
			}),
			(0, d.jsx)(`p`, { className: Pt.text, children: e.text }),
			e.status === `resolved` && (0, d.jsx)(A, { tone: `success`, children: `resolved` }),
		],
	});
}
var It = {
	root: `_root_d81sc_1`,
	item: `_item_d81sc_2`,
	button: `_button_d81sc_3`,
	number: `_number_d81sc_6`,
	body: `_body_d81sc_7`,
	title: `_title_d81sc_8`,
	description: `_description_d81sc_9`,
};
function Lt({ items: e, activeItemId: t, onItemSelect: n, className: r, ...i }) {
	return (0, d.jsx)(`ol`, {
		className: [It.root, r ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `CourseStepNav`,
		...i,
		children: e.map((e) => {
			let r = e.id === t;
			return (0, d.jsx)(
				`li`,
				{
					className: It.item,
					"data-active": r ? `true` : void 0,
					children: (0, d.jsxs)(`button`, {
						type: `button`,
						className: It.button,
						onClick: () => n?.(e.id),
						children: [
							(0, d.jsx)(`span`, { className: It.number, children: e.number }),
							(0, d.jsxs)(`span`, {
								className: It.body,
								children: [
									(0, d.jsx)(`span`, { className: It.title, children: e.title }),
									(0, d.jsx)(`span`, { className: It.description, children: e.description }),
								],
							}),
							e.duration && (0, d.jsx)(A, { children: e.duration }),
						],
					}),
				},
				e.id,
			);
		}),
	});
}
var Rt = { root: `_root_vz3kq_1`, cell: `_cell_vz3kq_2` };
function zt({ items: e, className: t, ...n }) {
	return (0, d.jsx)(`dl`, {
		className: [Rt.root, t ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `KeyValueStrip`,
		...n,
		children: e.map((e, t) =>
			(0, d.jsxs)(
				`div`,
				{
					className: Rt.cell,
					children: [
						(0, d.jsx)(`dt`, { children: e.key }),
						(0, d.jsx)(`dd`, { children: e.value }),
					],
				},
				t,
			),
		),
	});
}
var Bt = { root: `_root_yhnta_1`, item: `_item_yhnta_2`, marker: `_marker_yhnta_3` };
function Vt({ items: e, marker: t = `✔`, className: n, ...r }) {
	return (0, d.jsx)(`ul`, {
		className: [Bt.root, n ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `CheckList`,
		...r,
		children: e.map((e, n) => {
			let r = typeof e == `object` && e && `text` in e ? e.text : e,
				i = typeof e == `object` && e && `id` in e && e.id ? e.id : n;
			return (0, d.jsxs)(
				`li`,
				{
					className: Bt.item,
					children: [
						(0, d.jsx)(`span`, { className: Bt.marker, children: t }),
						(0, d.jsx)(`span`, { children: r }),
					],
				},
				i,
			);
		}),
	});
}
var Ht = {
	root: `_root_1ypbi_1`,
	item: `_item_1ypbi_2`,
	button: `_button_1ypbi_3`,
	index: `_index_1ypbi_7`,
	body: `_body_1ypbi_8`,
	title: `_title_1ypbi_9`,
	description: `_description_1ypbi_10`,
	meta: `_meta_1ypbi_11`,
};
function Ut({ items: e, activeItemId: t, onItemSelect: n, className: r, ...i }) {
	return (0, d.jsx)(`ol`, {
		className: [Ht.root, r ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `StepList`,
		...i,
		children: e.map((e) => {
			let r = e.id === t,
				i = (0, d.jsxs)(d.Fragment, {
					children: [
						(0, d.jsx)(`span`, { className: Ht.index, children: e.index }),
						(0, d.jsxs)(`span`, {
							className: Ht.body,
							children: [
								(0, d.jsx)(`span`, { className: Ht.title, children: e.title }),
								e.description &&
									(0, d.jsx)(`span`, { className: Ht.description, children: e.description }),
							],
						}),
						e.meta && (0, d.jsx)(`span`, { className: Ht.meta, children: e.meta }),
					],
				});
			return (0, d.jsx)(
				`li`,
				{
					className: Ht.item,
					"data-active": r ? `true` : void 0,
					children: n
						? (0, d.jsx)(`button`, {
								type: `button`,
								className: Ht.button,
								onClick: () => n(e.id),
								children: i,
							})
						: (0, d.jsx)(`div`, { className: Ht.button, children: i }),
				},
				e.id,
			);
		}),
	});
}
var Wt = {
	root: `_root_9v8jy_1`,
	avatar: `_avatar_9v8jy_2`,
	body: `_body_9v8jy_3`,
	name: `_name_9v8jy_4`,
	role: `_role_9v8jy_5`,
	bio: `_bio_9v8jy_6`,
};
function Gt({ name: e, subtitle: t, bio: n, avatar: r, className: i, ...a }) {
	let o = typeof e == `string` && e.length > 0 ? e[0] : `?`;
	return (0, d.jsxs)(`div`, {
		className: [Wt.root, i ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `PersonSummary`,
		...a,
		children: [
			(0, d.jsx)(`div`, { className: Wt.avatar, children: r ?? o }),
			(0, d.jsxs)(`div`, {
				className: Wt.body,
				children: [
					(0, d.jsx)(`div`, { className: Wt.name, children: e }),
					t && (0, d.jsx)(`div`, { className: Wt.role, children: t }),
					n && (0, d.jsx)(`div`, { className: Wt.bio, children: n }),
				],
			}),
		],
	});
}
var Kt = {
	root: `_root_zdnh5_1`,
	media: `_media_zdnh5_9`,
	bordered: `_bordered_zdnh5_13`,
	inset: `_inset_zdnh5_19`,
	supporting: `_supporting_zdnh5_25`,
	caption: `_caption_zdnh5_32`,
	legend: `_legend_zdnh5_36`,
};
function qt({
	as: e = `figure`,
	children: t,
	caption: n,
	legend: r,
	frame: i = `none`,
	className: a,
	...o
}) {
	return (0, d.jsxs)(e, {
		className: [Kt.root, Kt[i], a ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `FigureBlock`,
		...o,
		children: [
			(0, d.jsx)(`div`, { className: Kt.media, children: t }),
			(n || r) &&
				(0, d.jsxs)(`div`, {
					className: Kt.supporting,
					children: [
						n && (0, d.jsx)(A, { className: Kt.caption, children: n }),
						r && (0, d.jsx)(`div`, { className: Kt.legend, children: r }),
					],
				}),
		],
	});
}
var Jt = {
	root: `_root_u4dnv_1`,
	item: `_item_u4dnv_9`,
	index: `_index_u4dnv_17`,
	body: `_body_u4dnv_28`,
	title: `_title_u4dnv_35`,
	text: `_text_u4dnv_40`,
	meta: `_meta_u4dnv_45`,
};
function Yt(e, t) {
	return typeof e == `object` && e && !Array.isArray(e) && `text` in e
		? e
		: { text: e, index: String(t + 1).padStart(2, `0`) };
}
function Xt({ items: e, markerTone: t = `accent`, className: n, ...r }) {
	return (0, d.jsx)(`ol`, {
		className: [Jt.root, n ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `KeyPointList`,
		"data-marker-tone": t,
		...r,
		children: e.map((e, t) => {
			let n = Yt(e, t),
				r = n.index ?? String(t + 1).padStart(2, `0`);
			return (0, d.jsxs)(
				`li`,
				{
					className: Jt.item,
					children: [
						(0, d.jsx)(`span`, { className: Jt.index, children: r }),
						(0, d.jsxs)(`span`, {
							className: Jt.body,
							children: [
								n.title && (0, d.jsx)(`span`, { className: Jt.title, children: n.title }),
								(0, d.jsx)(`span`, { className: Jt.text, children: n.text }),
							],
						}),
						n.meta && (0, d.jsx)(`span`, { className: Jt.meta, children: n.meta }),
					],
				},
				n.id ?? t,
			);
		}),
	});
}
var Zt = {
	root: `_root_1y0w9_1`,
	section: `_section_1y0w9_9`,
	sectionLabel: `_sectionLabel_1y0w9_16`,
	items: `_items_1y0w9_24`,
	item: `_item_1y0w9_24`,
	icon: `_icon_1y0w9_67`,
	label: `_label_1y0w9_77`,
	badge: `_badge_1y0w9_83`,
};
function Qt({
	sections: e,
	activeItemId: t,
	onItemSelect: n,
	ariaLabel: r = `Course`,
	className: i,
	...a
}) {
	return (0, d.jsx)(`nav`, {
		className: [Zt.root, i ?? ``].filter(Boolean).join(` `),
		"aria-label": r,
		"data-rag-molecule": `SidebarNav`,
		...a,
		children: e.map((e) =>
			(0, d.jsxs)(
				`section`,
				{
					className: Zt.section,
					children: [
						(0, d.jsx)(`div`, { className: Zt.sectionLabel, children: e.label }),
						(0, d.jsx)(`div`, {
							className: Zt.items,
							children: e.items.map((e) => {
								let r = e.id === t;
								return (0, d.jsxs)(
									`button`,
									{
										type: `button`,
										className: Zt.item,
										"data-active": r ? `true` : void 0,
										disabled: e.disabled,
										"aria-current": r ? `page` : void 0,
										onClick: () => n?.(e.id),
										children: [
											(0, d.jsx)(`span`, {
												className: Zt.icon,
												"aria-hidden": `true`,
												children: e.icon,
											}),
											(0, d.jsx)(`span`, { className: Zt.label, children: e.label }),
											e.badge && (0, d.jsx)(`span`, { className: Zt.badge, children: e.badge }),
										],
									},
									e.id,
								);
							}),
						}),
					],
				},
				e.id,
			),
		),
	});
}
var j = {
	root: `_root_859oh_1`,
	h1: `_h1_859oh_6`,
	h2: `_h2_859oh_7`,
	h3: `_h3_859oh_8`,
	p: `_p_859oh_9`,
	ul: `_ul_859oh_10`,
	ol: `_ol_859oh_11`,
	blockquote: `_blockquote_859oh_12`,
	table: `_table_859oh_13`,
	codeBlock: `_codeBlock_859oh_14`,
	checkList: `_checkList_859oh_60`,
	check: `_check_859oh_60`,
	hr: `_hr_859oh_114`,
	inlineCode: `_inlineCode_859oh_120`,
	link: `_link_859oh_140`,
};
function $t(e, t) {
	let n = [],
		r = e,
		i = 0,
		a = /(\*\*([^*]+)\*\*|`([^`]+)`|\[([^\]]+)\]\(([^)]+)\))/,
		o;
	for (; (o = r.match(a)); )
		o.index && o.index > 0 && n.push(r.slice(0, o.index)),
			o[2] == null
				? o[3] == null
					? o[4] != null &&
						o[5] != null &&
						n.push(
							(0, d.jsx)(`a`, { href: o[5], className: j.link, children: o[4] }, `${t}-link-${i}`),
						)
					: n.push(
							(0, d.jsx)(`code`, { className: j.inlineCode, children: o[3] }, `${t}-code-${i}`),
						)
				: n.push((0, d.jsx)(`strong`, { children: o[2] }, `${t}-strong-${i}`)),
			(r = r.slice((o.index ?? 0) + o[0].length)),
			(i += 1);
	return r && n.push(r), n;
}
function en(e) {
	return e.every((e) => /^:?-{3,}:?$/.test(e.trim()));
}
function tn(e) {
	return e
		.split(`|`)
		.slice(1, -1)
		.map((e) => e.trim());
}
function nn({ source: e, className: t, ...n }) {
	let r = e
			.replace(
				/\r\n/g,
				`
`,
			)
			.split(`
`),
		i = [],
		a = 0,
		o = 0;
	for (; a < r.length; ) {
		let e = r[a] ?? ``,
			t = e.trim();
		if (!t) {
			a += 1;
			continue;
		}
		if (t.startsWith("```")) {
			let e = t.slice(3).trim(),
				n = [];
			for (a += 1; a < r.length && !(r[a] ?? ``).trim().startsWith("```"); )
				n.push(r[a] ?? ``), (a += 1);
			(a += 1),
				i.push(
					(0, d.jsx)(
						`pre`,
						{
							className: j.codeBlock,
							"data-language": e || void 0,
							children: (0, d.jsx)(`code`, {
								children: n.join(`
`),
							}),
						},
						o++,
					),
				);
			continue;
		}
		if (t === `---` || t === `***`) {
			i.push((0, d.jsx)(`hr`, { className: j.hr }, o++)), (a += 1);
			continue;
		}
		if (t.startsWith(`### `)) {
			i.push((0, d.jsx)(`h3`, { className: j.h3, children: $t(t.slice(4), `h3-${o}`) }, o++)),
				(a += 1);
			continue;
		}
		if (t.startsWith(`## `)) {
			i.push((0, d.jsx)(`h2`, { className: j.h2, children: $t(t.slice(3), `h2-${o}`) }, o++)),
				(a += 1);
			continue;
		}
		if (t.startsWith(`# `)) {
			i.push((0, d.jsx)(`h1`, { className: j.h1, children: $t(t.slice(2), `h1-${o}`) }, o++)),
				(a += 1);
			continue;
		}
		if (t.startsWith(`> `)) {
			let e = [];
			for (; a < r.length && (r[a] ?? ``).trim().startsWith(`> `); )
				e.push((r[a] ?? ``).trim().slice(2)), (a += 1);
			i.push(
				(0, d.jsx)(
					`blockquote`,
					{ className: j.blockquote, children: $t(e.join(` `), `quote-${o}`) },
					o++,
				),
			);
			continue;
		}
		if (t.startsWith(`|`)) {
			let e = [];
			for (; a < r.length && (r[a] ?? ``).trim().startsWith(`|`); )
				e.push(tn((r[a] ?? ``).trim())), (a += 1);
			let [t, ...n] = e.filter((e) => !en(e));
			t &&
				i.push(
					(0, d.jsxs)(
						`table`,
						{
							className: j.table,
							children: [
								(0, d.jsx)(`thead`, {
									children: (0, d.jsx)(`tr`, {
										children: t.map((e, t) =>
											(0, d.jsx)(`th`, { children: $t(e, `th-${o}-${t}`) }, t),
										),
									}),
								}),
								(0, d.jsx)(`tbody`, {
									children: n.map((e, t) =>
										(0, d.jsx)(
											`tr`,
											{
												children: e.map((e, n) =>
													(0, d.jsx)(`td`, { children: $t(e, `td-${o}-${t}-${n}`) }, n),
												),
											},
											t,
										),
									),
								}),
							],
						},
						o++,
					),
				);
			continue;
		}
		if (/^\s*-\s+\[[ xX]\]\s+/.test(e)) {
			let e = [];
			for (; a < r.length && /^\s*-\s+\[[ xX]\]\s+/.test(r[a] ?? ``); ) {
				let t = r[a] ?? ``;
				e.push({ checked: /\[[xX]\]/.test(t), text: t.replace(/^\s*-\s+\[[ xX]\]\s+/, ``) }),
					(a += 1);
			}
			i.push(
				(0, d.jsx)(
					`ul`,
					{
						className: j.checkList,
						children: e.map((e, t) =>
							(0, d.jsxs)(
								`li`,
								{
									children: [
										(0, d.jsx)(`span`, { className: j.check, children: e.checked ? `☑` : `☐` }),
										(0, d.jsx)(`span`, { children: $t(e.text, `check-${o}-${t}`) }),
									],
								},
								t,
							),
						),
					},
					o++,
				),
			);
			continue;
		}
		if (/^\s*-\s+/.test(e)) {
			let e = [];
			for (; a < r.length && /^\s*-\s+/.test(r[a] ?? ``); )
				e.push((r[a] ?? ``).replace(/^\s*-\s+/, ``)), (a += 1);
			i.push(
				(0, d.jsx)(
					`ul`,
					{
						className: j.ul,
						children: e.map((e, t) => (0, d.jsx)(`li`, { children: $t(e, `ul-${o}-${t}`) }, t)),
					},
					o++,
				),
			);
			continue;
		}
		if (/^\s*\d+\.\s+/.test(e)) {
			let e = [];
			for (; a < r.length && /^\s*\d+\.\s+/.test(r[a] ?? ``); )
				e.push((r[a] ?? ``).replace(/^\s*\d+\.\s+/, ``)), (a += 1);
			i.push(
				(0, d.jsx)(
					`ol`,
					{
						className: j.ol,
						children: e.map((e, t) => (0, d.jsx)(`li`, { children: $t(e, `ol-${o}-${t}`) }, t)),
					},
					o++,
				),
			);
			continue;
		}
		let n = [];
		for (; a < r.length; ) {
			let e = (r[a] ?? ``).trim();
			if (!e || /^(#{1,3}\s|>|```|---|\*\*\*|\||\s*-\s+|\s*\d+\.\s+)/.test(e)) break;
			n.push(e), (a += 1);
		}
		i.push((0, d.jsx)(`p`, { className: j.p, children: $t(n.join(` `), `p-${o}`) }, o++));
	}
	return (0, d.jsx)(`article`, {
		className: [j.root, t ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `MarkdownArticle`,
		...n,
		children: i,
	});
}
var rn = {
	root: `_root_iqxsi_1`,
	header: `_header_iqxsi_9`,
	description: `_description_iqxsi_15`,
	list: `_list_iqxsi_20`,
	empty: `_empty_iqxsi_26`,
	item: `_item_iqxsi_30`,
	icon: `_icon_iqxsi_66`,
	itemBody: `_itemBody_iqxsi_71`,
	itemTitle: `_itemTitle_iqxsi_78`,
	itemMeta: `_itemMeta_iqxsi_84`,
	itemDescription: `_itemDescription_iqxsi_90`,
	footer: `_footer_iqxsi_96`,
	downloadButton: `_downloadButton_iqxsi_103`,
};
function an(e) {
	let t = String(e ?? ``).toLowerCase();
	return t.includes(`pdf`)
		? `▤`
		: t.includes(`json`)
			? `{ }`
			: t.includes(`markdown`) || t.includes(`md`)
				? `¶`
				: `□`;
}
function on({
	title: e,
	description: t,
	items: n,
	selectedItemId: r,
	onItemSelect: i,
	onDownloadAll: a,
	downloadAllLabel: o = `⤓ Download all (.zip)`,
	emptyMessage: s = `No documents available.`,
	showItemDescriptions: c = !1,
	className: l,
	...u
}) {
	return (0, d.jsxs)(`div`, {
		className: [rn.root, l ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `DocumentListPanel`,
		...u,
		children: [
			(0, d.jsxs)(`header`, {
				className: rn.header,
				children: [
					(0, d.jsx)(_e, { as: `div`, size: `metric`, weight: `bold`, children: e }),
					t && (0, d.jsx)(A, { className: rn.description, children: t }),
				],
			}),
			(0, d.jsxs)(`div`, {
				className: rn.list,
				role: `listbox`,
				"aria-label": `Documents`,
				children: [
					n.length === 0 &&
						(0, d.jsx)(`div`, { className: rn.empty, children: (0, d.jsx)(A, { children: s }) }),
					n.map((e) => {
						let t = e.id === r;
						return (0, d.jsxs)(
							`button`,
							{
								type: `button`,
								className: rn.item,
								"data-active": t ? `true` : void 0,
								role: `option`,
								"aria-selected": t,
								disabled: e.disabled,
								onClick: () => i?.(e.id),
								children: [
									(0, d.jsx)(`span`, {
										className: rn.icon,
										"aria-hidden": `true`,
										children: e.icon ?? an(e.format),
									}),
									(0, d.jsxs)(`span`, {
										className: rn.itemBody,
										children: [
											(0, d.jsx)(`span`, { className: rn.itemTitle, children: e.title }),
											(0, d.jsx)(`span`, {
												className: rn.itemMeta,
												children: [e.format, e.size].filter(Boolean).join(` · `),
											}),
											c &&
												e.description &&
												(0, d.jsx)(`span`, {
													className: rn.itemDescription,
													children: e.description,
												}),
										],
									}),
								],
							},
							e.id,
						);
					}),
				],
			}),
			a &&
				(0, d.jsx)(`footer`, {
					className: rn.footer,
					children: (0, d.jsx)(f, {
						variant: `primary`,
						onClick: a,
						className: rn.downloadButton,
						children: o,
					}),
				}),
		],
	});
}
var sn = {
	root: `_root_17s2o_1`,
	meta: `_meta_17s2o_15`,
	actions: `_actions_17s2o_16`,
	badge: `_badge_17s2o_31`,
};
function cn({
	file: e,
	format: t,
	size: n,
	onDownload: r,
	downloadLabel: i = `⤓ Download`,
	rightSlot: a,
	className: o,
	...s
}) {
	return (0, d.jsxs)(`div`, {
		className: [sn.root, o ?? ``].filter(Boolean).join(` `),
		"data-rag-molecule": `DocumentPreviewToolbar`,
		...s,
		children: [
			(0, d.jsxs)(`div`, {
				className: sn.meta,
				children: [
					(0, d.jsx)(be, { children: e }),
					t &&
						(0, d.jsx)(`span`, {
							className: sn.badge,
							children: (0, d.jsx)(A, { tone: `inherit`, children: t }),
						}),
					n &&
						(0, d.jsx)(`span`, {
							className: sn.badge,
							children: (0, d.jsx)(A, { tone: `inherit`, children: n }),
						}),
				],
			}),
			(0, d.jsxs)(`div`, {
				className: sn.actions,
				children: [a, r && (0, d.jsx)(f, { size: `compact`, onClick: r, children: i })],
			}),
		],
	});
}
var ln = {
	root: `_root_88r47_1`,
	input: `_input_88r47_40`,
	icon: `_icon_88r47_54`,
	title: `_title_88r47_60`,
	description: `_description_88r47_64`,
};
function un(e) {
	return e ? Array.from(e) : [];
}
function dn({
	title: e = `Drop a file here`,
	description: t = `or choose a file`,
	icon: n,
	accept: r,
	multiple: i = !1,
	disabled: a = !1,
	active: o = !1,
	inputAriaLabel: s = `Choose file to upload`,
	onFilesSelected: c,
	className: l,
	onDragEnter: u,
	onDragOver: f,
	onDragLeave: p,
	onKeyDown: m,
	...h
}) {
	let g = (0, O.useRef)(null),
		[_, v] = (0, O.useState)(!1),
		y = o || _;
	function b() {
		a || g.current?.click();
	}
	function x(e) {
		!a && e.length > 0 && c?.(e);
	}
	function ee(e) {
		u?.(e), !(e.defaultPrevented || a) && (e.preventDefault(), v(!0));
	}
	function S(e) {
		f?.(e), !(e.defaultPrevented || a) && (e.preventDefault(), v(!0));
	}
	function C(e) {
		p?.(e), !a && (e.currentTarget.contains(e.relatedTarget) || v(!1));
	}
	function te(e) {
		e.preventDefault(), v(!1), x(un(e.dataTransfer.files));
	}
	function ne(e) {
		e.target !== g.current && b();
	}
	function re(e) {
		m?.(e),
			!(e.defaultPrevented || a) &&
				(e.key === `Enter` || e.key === ` `) &&
				(e.preventDefault(), b());
	}
	return (0, d.jsxs)(`div`, {
		className: [ln.root, l ?? ``].filter(Boolean).join(` `),
		role: `button`,
		tabIndex: a ? void 0 : 0,
		"aria-disabled": a || void 0,
		"data-active": y ? `true` : void 0,
		"data-disabled": a ? `true` : void 0,
		"data-rag-molecule": `FileDropZone`,
		onClick: ne,
		onDragEnter: ee,
		onDragOver: S,
		onDragLeave: C,
		onDrop: te,
		onKeyDown: re,
		...h,
		children: [
			(0, d.jsx)(`input`, {
				ref: g,
				className: ln.input,
				type: `file`,
				accept: r,
				multiple: i,
				disabled: a,
				"aria-label": s,
				tabIndex: -1,
				onChange: (e) => {
					x(un(e.currentTarget.files)), (e.currentTarget.value = ``);
				},
			}),
			(0, d.jsx)(`div`, {
				className: ln.icon,
				"aria-hidden": `true`,
				children: n ?? (0, d.jsx)(E, { title: `` }),
			}),
			(0, d.jsx)(_e, { as: `div`, size: `body`, weight: `bold`, className: ln.title, children: e }),
			t && (0, d.jsx)(A, { className: ln.description, children: t }),
		],
	});
}
var fn = { viewport: `_viewport_1j377_1` },
	pn = [`strip`, `stack`, `budget`, `treemap`];
function mn({ snapshot: e, initialView: t = `strip`, selectedPartId: n, className: r, ...i }) {
	let [a, o] = (0, O.useState)(t),
		s = n ?? e.selectedPartId;
	return (0, d.jsx)(ke, {
		title: e.title,
		actions: (0, d.jsx)(Le, {
			gap: `xs`,
			children: pn.map((e) =>
				(0, d.jsx)(
					f,
					{
						size: `compact`,
						selected: a === e,
						"aria-pressed": a === e,
						onClick: () => o(e),
						children: e,
					},
					e,
				),
			),
		}),
		className: r,
		"data-rag-organism": `ContextDiagramPanel`,
		...i,
		children: (0, d.jsxs)(Ne, {
			gap: `sm`,
			children: [
				e.subtitle && (0, d.jsx)(A, { children: e.subtitle }),
				(0, d.jsxs)(`div`, {
					className: fn.viewport,
					children: [
						a === `strip` && (0, d.jsx)(vt, { snapshot: e, selectedPartId: s }),
						a === `stack` && (0, d.jsx)(xt, { snapshot: e, selectedPartId: s }),
						a === `budget` && (0, d.jsx)(ht, { snapshot: e, selectedPartId: s }),
						a === `treemap` && (0, d.jsx)(wt, { snapshot: e, selectedPartId: s }),
					],
				}),
				(0, d.jsx)(dt, {
					compact: !0,
					kinds: [
						`system`,
						`context`,
						`summary`,
						`result`,
						`generated`,
						`evicted`,
						`active`,
						`empty`,
					],
					selectedKind: s ? e.parts.find((e) => e.id === s)?.kind : void 0,
				}),
			],
		}),
	});
}
var hn = { root: `_root_1szio_1`, stream: `_stream_1szio_9` };
function gn(e) {
	return e.reduce((e, t) => e + (t.tokens ?? 0), 0);
}
function _n({
	title: e = `Transcript`,
	subtitle: t,
	messages: n,
	annotations: r = [],
	selectedAnnotationId: i,
	onAnnotationSelect: a,
	showAnnotationChips: o = !0,
	className: s,
	...c
}) {
	return (0, d.jsxs)(`section`, {
		className: [hn.root, s ?? ``].filter(Boolean).join(` `),
		"data-rag-organism": `TranscriptReaderPanel`,
		...c,
		children: [
			(0, d.jsx)(jt, {
				title: e,
				subtitle: t,
				messageCount: n.length,
				annotationCount: r.length,
				tokenTotal: gn(n),
			}),
			(0, d.jsx)(`div`, {
				className: hn.stream,
				children: n.map((e) =>
					(0, d.jsx)(
						kt,
						{
							message: e,
							annotations: r,
							selectedAnnotationId: i,
							onAnnotationSelect: a,
							showAnnotationChips: o,
						},
						e.id,
					),
				),
			}),
		],
	});
}
var vn = {
	root: `_root_1jf6m_1`,
	header: `_header_1jf6m_10`,
	description: `_description_1jf6m_17`,
	list: `_list_1jf6m_22`,
	button: `_button_1jf6m_32`,
};
function yn({
	title: e = `Context-engineering notes`,
	description: t = `What each message costs the window. Click a note or a message to focus it.`,
	annotations: n,
	selectedAnnotationId: r,
	onAnnotationSelect: i,
	className: a,
	...o
}) {
	return (0, d.jsxs)(`aside`, {
		className: [vn.root, a ?? ``].filter(Boolean).join(` `),
		"data-rag-organism": `AnnotationRailPanel`,
		...o,
		children: [
			(0, d.jsxs)(`header`, {
				className: vn.header,
				children: [
					(0, d.jsx)(_e, { as: `div`, size: `label`, children: e }),
					(0, d.jsx)(A, { className: vn.description, children: t }),
				],
			}),
			(0, d.jsx)(`div`, {
				className: vn.list,
				children: n.map((e, t) =>
					(0, d.jsx)(
						`button`,
						{
							type: `button`,
							className: vn.button,
							onClick: () => i?.(e.id),
							children: (0, d.jsx)(Nt, { annotation: e, selected: e.id === r, index: t + 1 }),
						},
						e.id,
					),
				),
			}),
		],
	});
}
var bn = { root: `_root_1q01w_1`, withNotes: `_withNotes_1q01w_8`, noNotes: `_noNotes_1q01w_13` };
function xn({
	title: e,
	subtitle: t,
	messages: n,
	annotations: r = [],
	selectedAnnotationId: i,
	onAnnotationSelect: a,
	showNotes: o = r.length > 0,
	className: s,
	...c
}) {
	let l = o && r.length > 0;
	return (0, d.jsxs)(`div`, {
		className: [bn.root, l ? bn.withNotes : bn.noNotes, s ?? ``].filter(Boolean).join(` `),
		"data-rag-organism": `TranscriptWorkspacePanel`,
		...c,
		children: [
			(0, d.jsx)(_n, {
				title: e,
				subtitle: t,
				messages: n,
				annotations: l ? r : [],
				selectedAnnotationId: i,
				onAnnotationSelect: a,
				showAnnotationChips: l,
			}),
			l && (0, d.jsx)(yn, { annotations: r, selectedAnnotationId: i, onAnnotationSelect: a }),
		],
	});
}
function Sn({
	title: e = `Comments`,
	comments: t,
	selectedCommentId: n,
	onCommentSelect: r,
	...i
}) {
	return (0, d.jsx)(ke, {
		title: e,
		fill: !0,
		"data-rag-organism": `AnchoredCommentRail`,
		...i,
		children: (0, d.jsxs)(Ne, {
			gap: `sm`,
			children: [
				(0, d.jsxs)(A, { children: [t.length, ` anchored comment`, t.length === 1 ? `` : `s`] }),
				t.map((e, t) =>
					(0, d.jsx)(
						`button`,
						{
							type: `button`,
							style: {
								appearance: `none`,
								border: 0,
								background: `transparent`,
								padding: 0,
								textAlign: `left`,
								cursor: `pointer`,
							},
							onClick: () => r?.(e.id),
							children: (0, d.jsx)(Ft, { comment: e, index: t + 1, selected: e.id === n }),
						},
						e.id,
					),
				),
			],
		}),
	});
}
var Cn = {
	root: `_root_1io7b_1`,
	hero: `_hero_1io7b_8`,
	eyebrow: `_eyebrow_1io7b_12`,
	title: `_title_1io7b_20`,
	lead: `_lead_1io7b_28`,
	ctaRow: `_ctaRow_1io7b_35`,
	metaStrip: `_metaStrip_1io7b_42`,
	instructorBlock: `_instructorBlock_1io7b_47`,
	sectionLabel: `_sectionLabel_1io7b_51`,
	blurb: `_blurb_1io7b_59`,
};
function wn({ course: e, activeAgendaItemId: t, onAgendaItemSelect: n, className: r, ...i }) {
	let a = e.agenda.map((e) => ({
		id: e.id,
		index: e.number,
		title: e.title,
		description: e.description,
		meta: e.duration,
	}));
	return (0, d.jsxs)(`article`, {
		className: [Cn.root, r ?? ``].filter(Boolean).join(` `),
		"data-rag-organism": `CourseLessonPanel`,
		...i,
		children: [
			(0, d.jsxs)(Xe, {
				density: `spacious`,
				divider: `bottom`,
				className: Cn.hero,
				children: [
					e.kicker && (0, d.jsx)(`div`, { className: Cn.eyebrow, children: e.kicker }),
					(0, d.jsx)(`h1`, { className: Cn.title, children: e.title }),
					(0, d.jsx)(`p`, { className: Cn.lead, children: e.tagline }),
					(0, d.jsxs)(`div`, {
						className: Cn.ctaRow,
						children: [
							(0, d.jsx)(f, { variant: `primary`, children: `Reserve a seat` }),
							(0, d.jsx)(f, { children: `Preview the deck` }),
							e.price && (0, d.jsx)(A, { children: e.price }),
						],
					}),
					(0, d.jsx)(zt, {
						className: Cn.metaStrip,
						items: [
							{ key: `When`, value: e.when ?? `TBD` },
							{ key: `Where`, value: e.where ?? `TBD` },
							{ key: `Format`, value: e.format ?? `TBD` },
						],
					}),
				],
			}),
			(0, d.jsx)(Qe, {
				ratio: `course`,
				divider: !0,
				left: (0, d.jsxs)(Xe, {
					label: `What you'll leave with`,
					density: `spacious`,
					children: [
						(0, d.jsx)(Vt, { items: e.outcomes }),
						e.instructor &&
							(0, d.jsxs)(`div`, {
								className: Cn.instructorBlock,
								children: [
									(0, d.jsx)(`div`, { className: Cn.sectionLabel, children: `Instructor` }),
									(0, d.jsx)(Gt, {
										name: e.instructor.name,
										subtitle: e.instructor.role,
										bio: e.instructor.bio,
									}),
								],
							}),
					],
				}),
				right: (0, d.jsx)(Xe, {
					label: `Agenda · 2 hours`,
					density: `spacious`,
					children: (0, d.jsx)(Ut, { items: a, activeItemId: t, onItemSelect: n }),
				}),
			}),
			e.blurb &&
				(0, d.jsx)(Xe, {
					density: `spacious`,
					divider: `top`,
					children: (0, d.jsx)(`p`, { className: Cn.blurb, children: e.blurb }),
				}),
		],
	});
}
function Tn(e, t) {
	return e.view === `strip`
		? (0, d.jsx)(vt, { snapshot: t })
		: e.view === `stack`
			? (0, d.jsx)(xt, { snapshot: t })
			: e.view === `budget`
				? (0, d.jsx)(ht, { snapshot: t })
				: (0, d.jsx)(wt, { snapshot: t });
}
function En({
	slide: e,
	snapshot: t,
	index: n,
	total: r,
	visualSide: i = `left`,
	onPrevious: a,
	onNext: o,
	className: s,
	...c
}) {
	let l =
			n != null && r != null
				? `${String(n + 1).padStart(2, `0`)} / ${String(r).padStart(2, `0`)}`
				: void 0,
		u = (0, d.jsx)(qt, {
			caption: `${t.title} — ${t.limit.toLocaleString()} token window`,
			legend: (0, d.jsx)(dt, { compact: !0 }),
			children: Tn(e, t),
		}),
		p = (0, d.jsx)(Xt, { items: e.notes }),
		m =
			a || o
				? (0, d.jsxs)(Le, {
						justify: `between`,
						children: [
							(0, d.jsx)(f, { size: `compact`, onClick: a, disabled: !a, children: `◂ Prev` }),
							(0, d.jsx)(f, { size: `compact`, onClick: o, disabled: !o, children: `Next ▸` }),
						],
					})
				: void 0;
	return (0, d.jsx)(et, {
		eyebrow: e.kicker,
		counter: l,
		title: e.title,
		primary: u,
		secondary: p,
		primarySide: i,
		ratio: `primaryWide`,
		footer: m,
		className: s,
		"data-rag-organism": `CourseSlidePanel`,
		...c,
	});
}
var Dn = { header: `_header_i94hr_1`, title: `_title_i94hr_7`, content: `_content_i94hr_13` };
function On({
	sections: e,
	activeItemId: t,
	onNavigate: n,
	title: r = `Context Window Engineering`,
	subtitle: i = `Course studio`,
	sidebarFooter: a,
	className: o,
	children: s,
	...c
}) {
	return (0, d.jsx)(rt, {
		className: o,
		sidebarWidth: 188,
		header: (0, d.jsxs)(`div`, {
			className: Dn.header,
			children: [
				(0, d.jsx)(`div`, { className: Dn.title, children: r }),
				i && (0, d.jsx)(A, { children: i }),
			],
		}),
		footer: a,
		sidebar: (0, d.jsx)(Qt, { sections: e, activeItemId: t, onItemSelect: n }),
		"data-rag-organism": `CourseStudioShell`,
		...c,
		children: (0, d.jsx)(`div`, { className: Dn.content, children: s }),
	});
}
var kn = {
	root: `_root_1a4mz_1`,
	sidebar: `_sidebar_1a4mz_10`,
	preview: `_preview_1a4mz_17`,
	articleWrap: `_articleWrap_1a4mz_24`,
	placeholder: `_placeholder_1a4mz_30`,
};
function An(e, t) {
	return e.find((e) => e.id === t) ?? e[0];
}
function jn({
	intro: e,
	documents: t,
	selectedDocumentId: n,
	onDocumentSelect: r,
	onDownload: i,
	onDownloadAll: a,
	title: o = `Handout`,
	emptyMessage: s = `Select a document to preview it.`,
	className: c,
	...l
}) {
	let u = An(t, n);
	return (0, d.jsxs)(`div`, {
		className: [kn.root, c ?? ``].filter(Boolean).join(` `),
		"data-rag-organism": `HandoutDocumentShell`,
		...l,
		children: [
			(0, d.jsx)(`aside`, {
				className: kn.sidebar,
				children: (0, d.jsx)(on, {
					title: o,
					description: e,
					items: t.map((e) => ({
						id: e.id,
						title: e.title,
						format: e.format,
						size: e.size,
						description: e.description,
					})),
					selectedItemId: u?.id,
					onItemSelect: r,
					onDownloadAll: a,
				}),
			}),
			(0, d.jsx)(`section`, {
				className: kn.preview,
				children: u
					? (0, d.jsxs)(d.Fragment, {
							children: [
								(0, d.jsx)(cn, {
									file: u.file,
									format: u.format,
									size: u.size,
									onDownload: i ? () => i(u.id) : void 0,
								}),
								(0, d.jsx)(`div`, {
									className: kn.articleWrap,
									children: u.format.toLowerCase().includes(`markdown`)
										? (0, d.jsx)(nn, { source: u.body })
										: (0, d.jsx)(`div`, {
												className: kn.placeholder,
												children: (0, d.jsx)(A, { children: u.description }),
											}),
								}),
							],
						})
					: (0, d.jsx)(`div`, {
							className: kn.placeholder,
							children: (0, d.jsx)(A, { children: s }),
						}),
			}),
		],
	});
}
function Mn({
	title: e = `Drop a .json file here`,
	description: t = `or paste below · max 200k tokens`,
	accept: n = `application/json,.json`,
	disabled: r = !1,
	active: i = !1,
	onFilesSelected: a,
	...o
}) {
	return (0, d.jsx)(dn, {
		title: e,
		description: t,
		accept: n,
		disabled: r,
		active: i,
		onFilesSelected: a,
		inputAriaLabel: `Choose context-window JSON file`,
		"data-rag-organism": `ContextUploadDropArea`,
		...o,
	});
}
var Nn = c();
function Pn(e, t = {}, n) {
	if (n) {
		n(e, t);
		return;
	}
	if (e.kind === `copy`) {
		let n = e.value ?? (e.field && t.row ? String(t.row[e.field] ?? ``) : ``);
		n && navigator.clipboard?.writeText(n).catch(() => {});
		return;
	}
	if (e.kind === `event`) {
		window.dispatchEvent(new CustomEvent(e.event, { detail: { ...(e.detail ?? {}), context: t } }));
		return;
	}
	if (e.kind === `navigate`) {
		let n = In(e.to, t);
		window.history.pushState(e.params ?? {}, ``, n),
			window.dispatchEvent(new PopStateEvent(`popstate`));
		return;
	}
	e.kind === `server` &&
		fetch(`/api/widget/actions/${encodeURIComponent(e.name)}`, {
			method: `POST`,
			headers: { "Content-Type": `application/json` },
			body: JSON.stringify({ payload: e.payload ?? {}, context: t }),
		});
}
function Fn(e, t, n) {
	if (e) return () => Pn(e, t, n);
}
function In(e, t) {
	return e.replace(/\$\{([^}]+)\}|\$([A-Za-z0-9_.-]+)/g, (e, n, r) => {
		let i = Ln(n ?? r ?? ``, t);
		return encodeURIComponent(String(i ?? ``));
	});
}
function Ln(e, t) {
	let n = e.split(`.`).filter(Boolean),
		r = t;
	for (let e of n) {
		if (!r || typeof r != `object`) return;
		r = r[e];
	}
	return r;
}
function Rn(e) {
	if (!e || typeof e != `object`) return !1;
	let t = e.kind;
	return t === `text` || t === `element` || t === `component`;
}
function zn(e, t, n) {
	switch (e.kind) {
		case `field`:
			return Un(Hn(t, e.field), e.fallback);
		case `number`:
			return Wn(Hn(t, e.field), e);
		case `status`: {
			let n = Un(Hn(t, e.field), e.fallback ?? `pending`);
			return (0, d.jsx)(Ce, { status: n, icon: e.icon, children: n });
		}
		case `caption`:
			return (0, d.jsx)(A, { tone: e.tone, children: Un(Hn(t, e.field), e.fallback) });
		case `template`:
			return Gn(e.template, t);
		case `link`: {
			let n = Un(Hn(t, e.hrefField), `#`),
				r = Un(Hn(t, e.labelField), e.fallbackLabel ?? n);
			return (0, d.jsx)(`a`, {
				href: n,
				target: e.target,
				rel: e.target === `_blank` ? `noreferrer` : void 0,
				children: r,
			});
		}
		case `constant`:
			return Bn(e.value, n);
		default:
			return (0, d.jsx)(A, { tone: `danger`, children: `Unsupported cell` });
	}
}
function Bn(e, t) {
	return e == null || e === !1 ? null : Rn(e) ? t(e) : String(e);
}
function Vn(e, t) {
	return typeof t == `string`
		? Un(Hn(e, t), ``)
		: `field` in t
			? Un(Hn(e, t.field), ``)
			: Gn(t.template, e);
}
function Hn(e, t) {
	let n = t.split(`.`).filter(Boolean),
		r = e;
	for (let e of n) {
		if (!r || typeof r != `object`) return;
		r = r[e];
	}
	return r;
}
function Un(e, t = ``) {
	return e == null ? t : String(e);
}
function Wn(e, t) {
	let n = typeof e == `number` ? e : Number(e);
	return Number.isFinite(n)
		? t.format === `fixed`
			? n.toFixed(t.digits ?? 2)
			: Math.round(n).toLocaleString()
		: (t.fallback ?? ``);
}
function Gn(e, t) {
	return e.replace(/\$\{([^}]+)\}|\$([A-Za-z0-9_.-]+)/g, (e, n, r) => Un(Hn(t, n ?? r ?? ``), ``));
}
function Kn({ node: e, onAction: t }) {
	return (0, d.jsx)(d.Fragment, { children: qn(e, t) });
}
function qn(e, t) {
	return e.kind === `text` ? e.text : e.kind === `element` ? Xn(e, t) : Zn(e, t);
}
function Jn(e, t) {
	return (e ?? []).map((e, n) => (0, d.jsx)(Kn, { node: e, onAction: t }, n));
}
function M(e, t) {
	return Bn(e, (e) => qn(e, t));
}
function Yn(e, t) {
	return e ? qn(e, t) : void 0;
}
function Xn(e, t) {
	let n = { ...(e.attrs ?? {}) };
	return (
		n.style && typeof n.style == `object` && (n.style = n.style),
		(0, O.createElement)(e.tag, n, ...Jn(e.children, t))
	);
}
function Zn(e, t) {
	switch (e.type) {
		case `AppShell`:
			return Qn(e, t);
		case `AppNav`:
			return $n(e, t);
		case `Button`:
			return er(e, t);
		case `Text`:
			return tr(e, t);
		case `CodeText`:
			return nr(e, t);
		case `Divider`:
			return rr(e);
		case `ContextKindSwatch`:
			return ir(e);
		case `ContextStudioNavIcon`:
			return ar(e);
		case `AnnotationBadge`:
			return or(e);
		case `TranscriptRoleBadge`:
			return sr(e);
		case `ContextLegend`:
			return cr(e);
		case `ContextBudgetBar`:
			return lr(e);
		case `ContextStripDiagram`:
			return ur(e);
		case `ContextStackDiagram`:
			return dr(e);
		case `ContextTreemap`:
			return fr(e);
		case `ContextDiagramPanel`:
			return pr(e);
		case `TranscriptSessionHeader`:
			return mr(e, t);
		case `TranscriptMessageCard`:
			return hr(e, t);
		case `AnnotationNoteCard`:
			return gr(e);
		case `AnnotationRailPanel`:
			return _r(e, t);
		case `TranscriptReaderPanel`:
			return vr(e, t);
		case `TranscriptWorkspacePanel`:
			return yr(e, t);
		case `AnchoredCommentCard`:
			return br(e, t);
		case `AnchoredCommentRail`:
			return xr(e, t);
		case `KeyValueStrip`:
			return Sr(e, t);
		case `CheckList`:
			return Cr(e, t);
		case `StepList`:
			return wr(e, t);
		case `PersonSummary`:
			return Tr(e, t);
		case `FigureBlock`:
			return Er(e, t);
		case `KeyPointList`:
			return Dr(e, t);
		case `SidebarNav`:
			return Or(e, t);
		case `CourseStepNav`:
			return kr(e, t);
		case `MarkdownArticle`:
			return Ar(e);
		case `DocumentListPanel`:
			return jr(e, t);
		case `DocumentPreviewToolbar`:
			return Mr(e, t);
		case `CourseLessonPanel`:
			return Nr(e, t);
		case `CourseSlidePanel`:
			return Pr(e, t);
		case `CourseStudioShell`:
			return Fr(e, t);
		case `HandoutDocumentShell`:
			return Ir(e, t);
		case `ContextUploadDropArea`:
			return Lr(e, t);
		case `Caption`:
			return Ur(e, t);
		case `DashboardGrid`:
			return Wr(e, t);
		case `DataTable`:
			return Gr(e, t);
		case `FormRow`:
			return Kr(e, t);
		case `Inline`:
			return qr(e, t);
		case `MetadataGrid`:
			return Jr(e, t);
		case `Panel`:
			return Yr(e, t);
		case `ScrollRegion`:
			return Xr(e, t);
		case `SectionBlock`:
			return Zr(e, t);
		case `SplitPane`:
			return Qr(e, t);
		case `SidebarShell`:
			return $r(e, t);
		case `SlideShell`:
			return ei(e, t);
		case `SelectInput`:
			return ti(e, t);
		case `Stack`:
			return ni(e, t);
		case `StatusText`:
			return ri(e, t);
		case `TabList`:
			return ii(e, t);
		case `TextInput`:
			return ai(e);
		default:
			return (0, d.jsxs)(m, { children: [`Unknown widget: `, e.type] });
	}
}
function Qn(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(He, {
		className: n.className,
		header: n.header ? qn(n.header, t) : void 0,
		sidebar: n.sidebar ? qn(n.sidebar, t) : void 0,
		children: Jn(e.children, t),
	});
}
function $n(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(lt, {
		brand: M(n.brand, t),
		items: n.items.map((e) => ({ id: e.id, label: M(e.label, t) })),
		activeItemId: n.activeItemId,
		onItemSelect: (e) => {
			let r = n.items.find((t) => t.id === e);
			r?.action && oi(r.action, { value: e, componentType: `AppNav` }, t);
		},
	});
}
function er(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(f, {
		className: n.className,
		variant: n.variant,
		size: n.size,
		selected: n.selected,
		disabled: n.disabled,
		type: n.type ?? `button`,
		onClick: Fn(n.action, { componentType: `Button` }, t),
		children: Jn(e.children, t),
	});
}
function tr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(_e, {
		className: n.className,
		as: n.as,
		size: n.size,
		tone: n.tone,
		weight: n.weight,
		align: n.align,
		truncate: n.truncate,
		children: Jn(e.children, t),
	});
}
function nr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(be, {
		className: n.className,
		as: n.as,
		tone: n.tone,
		display: n.display,
		copyable: n.copyable,
		children: Jn(e.children, t),
	});
}
function rr(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(Te, { className: t.className, orientation: t.orientation });
}
function ir(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(oe, {
		className: t.className,
		kind: t.kind,
		mode: t.mode,
		size: t.size,
		selected: t.selected,
	});
}
function ar(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(fe, { className: t.className, id: t.id, title: t.title });
}
function or(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(ue, {
		className: t.className,
		kind: t.kind,
		label: t.label,
		selected: t.selected,
	});
}
function sr(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(ce, { className: t.className, role: t.role, name: t.name });
}
function cr(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(dt, {
		className: t.className,
		kinds: t.kinds,
		mode: t.mode,
		compact: t.compact,
		selectedKind: t.selectedKind,
	});
}
function lr(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(ht, {
		className: t.className,
		snapshot: t.snapshot,
		mode: t.mode,
		showLegend: t.showLegend,
		selectedPartId: t.selectedPartId,
	});
}
function ur(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(vt, {
		className: t.className,
		snapshot: t.snapshot,
		mode: t.mode,
		selectedPartId: t.selectedPartId,
		showLabels: t.showLabels,
	});
}
function dr(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(xt, {
		className: t.className,
		snapshot: t.snapshot,
		selectedPartId: t.selectedPartId,
	});
}
function fr(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(wt, {
		className: t.className,
		snapshot: t.snapshot,
		selectedPartId: t.selectedPartId,
	});
}
function pr(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(mn, {
		className: t.className,
		snapshot: t.snapshot,
		initialView: t.initialView,
		selectedPartId: t.selectedPartId,
	});
}
function mr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(jt, {
		className: n.className,
		title: M(n.title, t),
		subtitle: M(n.subtitle, t),
		messageCount: n.messageCount,
		annotationCount: n.annotationCount,
		tokenTotal: n.tokenTotal,
		rightSlot: Yn(n.rightSlot, t),
	});
}
function hr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(kt, {
		className: n.className,
		message: n.message,
		annotations: n.annotations,
		selectedAnnotationId: n.selectedAnnotationId,
		showAnnotationChips: n.showAnnotationChips,
		onAnnotationSelect: Hr(`TranscriptMessageCard`, n.onAnnotationSelectAction, t),
	});
}
function gr(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(Nt, {
		className: t.className,
		annotation: t.annotation,
		selected: t.selected,
		index: t.index,
	});
}
function _r(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(yn, {
		className: n.className,
		title: n.title,
		description: n.description,
		annotations: n.annotations,
		selectedAnnotationId: n.selectedAnnotationId,
		onAnnotationSelect: Hr(`AnnotationRailPanel`, n.onAnnotationSelectAction, t),
	});
}
function vr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(_n, {
		className: n.className,
		title: n.title,
		subtitle: n.subtitle,
		messages: n.messages,
		annotations: n.annotations,
		selectedAnnotationId: n.selectedAnnotationId,
		showAnnotationChips: n.showAnnotationChips,
		onAnnotationSelect: Hr(`TranscriptReaderPanel`, n.onAnnotationSelectAction, t),
	});
}
function yr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(xn, {
		className: n.className,
		title: n.title,
		subtitle: n.subtitle,
		messages: n.messages,
		annotations: n.annotations,
		selectedAnnotationId: n.selectedAnnotationId,
		showNotes: n.showNotes,
		onAnnotationSelect: Hr(`TranscriptWorkspacePanel`, n.onAnnotationSelectAction, t),
	});
}
function br(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Ft, {
		className: n.className,
		comment: n.comment,
		index: n.index,
		selected: n.selected,
		compact: n.compact,
		onDismiss: n.onDismissAction
			? () =>
					oi(
						n.onDismissAction,
						{ commentId: n.comment.id, value: n.comment.id, componentType: `AnchoredCommentCard` },
						t,
					)
			: void 0,
	});
}
function xr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Sn, {
		className: n.className,
		title: n.title,
		comments: n.comments,
		selectedCommentId: n.selectedCommentId,
		onCommentSelect: n.onCommentSelectAction
			? (e) =>
					oi(
						n.onCommentSelectAction,
						{ commentId: e, value: e, componentType: `AnchoredCommentRail` },
						t,
					)
			: void 0,
	});
}
function Sr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(zt, {
		className: n.className,
		items: n.items.map((e) => ({ key: M(e.key, t), value: M(e.value, t) })),
	});
}
function Cr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Vt, {
		className: n.className,
		marker: M(n.marker, t),
		items: n.items.map((e) => (zr(e) ? { ...e, text: M(e.text, t) } : M(e, t))),
	});
}
function wr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Ut, {
		className: n.className,
		items: n.items.map((e) => ({
			...e,
			index: M(e.index, t),
			title: M(e.title, t),
			description: M(e.description, t),
			meta: M(e.meta, t),
		})),
		activeItemId: n.activeItemId,
		onItemSelect: Vr(`StepList`, n.onItemSelectAction, `itemId`, t),
	});
}
function Tr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Gt, {
		className: n.className,
		name: M(n.name, t),
		subtitle: M(n.subtitle, t),
		bio: M(n.bio, t),
		avatar: M(n.avatar, t),
	});
}
function Er(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(qt, {
		className: n.className,
		as: n.as,
		caption: M(n.caption, t),
		legend: Yn(n.legend, t),
		frame: n.frame,
		children: Jn(e.children, t),
	});
}
function Dr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Xt, {
		className: n.className,
		markerTone: n.markerTone,
		items: n.items.map((e) =>
			Br(e)
				? {
						...e,
						index: M(e.index, t),
						title: M(e.title, t),
						text: M(e.text, t),
						meta: M(e.meta, t),
					}
				: M(e, t),
		),
	});
}
function Or(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Qt, {
		className: n.className,
		sections: n.sections.map((e) => ({
			...e,
			label: M(e.label, t),
			items: e.items.map((e) => ({
				...e,
				label: M(e.label, t),
				icon: M(e.icon, t),
				badge: M(e.badge, t),
			})),
		})),
		activeItemId: n.activeItemId,
		ariaLabel: n.ariaLabel,
		onItemSelect: Vr(`SidebarNav`, n.onItemSelectAction, `itemId`, t),
	});
}
function kr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Lt, {
		className: n.className,
		items: n.items,
		activeItemId: n.activeItemId,
		onItemSelect: Vr(`CourseStepNav`, n.onItemSelectAction, `itemId`, t),
	});
}
function Ar(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(nn, { className: t.className, source: t.source });
}
function jr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(on, {
		className: n.className,
		title: M(n.title, t),
		description: M(n.description, t),
		items: n.items.map((e) => ({
			...e,
			title: M(e.title, t),
			format: M(e.format, t),
			size: M(e.size, t),
			description: M(e.description, t),
			icon: M(e.icon, t),
		})),
		selectedItemId: n.selectedItemId,
		onItemSelect: Vr(`DocumentListPanel`, n.onItemSelectAction, `itemId`, t),
		onDownloadAll: n.onDownloadAllAction
			? () => oi(n.onDownloadAllAction, { componentType: `DocumentListPanel` }, t)
			: void 0,
		downloadAllLabel: M(n.downloadAllLabel, t),
		emptyMessage: M(n.emptyMessage, t),
		showItemDescriptions: n.showItemDescriptions,
	});
}
function Mr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(cn, {
		className: n.className,
		file: M(n.file, t),
		format: M(n.format, t),
		size: M(n.size, t),
		onDownload: n.onDownloadAction
			? () => oi(n.onDownloadAction, { componentType: `DocumentPreviewToolbar` }, t)
			: void 0,
		downloadLabel: M(n.downloadLabel, t),
		rightSlot: Yn(n.rightSlot, t),
	});
}
function Nr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(wn, {
		className: n.className,
		course: n.course,
		activeAgendaItemId: n.activeAgendaItemId,
		onAgendaItemSelect: Vr(`CourseLessonPanel`, n.onAgendaItemSelectAction, `agendaItemId`, t),
	});
}
function Pr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(En, {
		className: n.className,
		slide: n.slide,
		snapshot: n.snapshot,
		index: n.index,
		total: n.total,
		visualSide: n.visualSide,
		onPrevious: n.onPreviousAction
			? () => oi(n.onPreviousAction, { componentType: `CourseSlidePanel`, value: `previous` }, t)
			: void 0,
		onNext: n.onNextAction
			? () => oi(n.onNextAction, { componentType: `CourseSlidePanel`, value: `next` }, t)
			: void 0,
	});
}
function Fr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(On, {
		className: n.className,
		sections: n.sections.map((e) => ({
			...e,
			label: M(e.label, t),
			items: e.items.map((e) => ({
				...e,
				label: M(e.label, t),
				icon: M(e.icon, t),
				badge: M(e.badge, t),
			})),
		})),
		activeItemId: n.activeItemId,
		onNavigate: Vr(`CourseStudioShell`, n.onNavigateAction, `itemId`, t),
		title: M(n.title, t),
		subtitle: M(n.subtitle, t),
		sidebarFooter: Yn(n.sidebarFooter, t),
		children: Jn(e.children, t),
	});
}
function Ir(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(jn, {
		className: n.className,
		intro: M(n.intro, t),
		documents: n.documents,
		selectedDocumentId: n.selectedDocumentId,
		onDocumentSelect: Vr(`HandoutDocumentShell`, n.onDocumentSelectAction, `documentId`, t),
		onDownload: n.onDownloadAction
			? (e) =>
					oi(
						n.onDownloadAction,
						{ documentId: e, value: e, componentType: `HandoutDocumentShell` },
						t,
					)
			: void 0,
		onDownloadAll: n.onDownloadAllAction
			? () => oi(n.onDownloadAllAction, { componentType: `HandoutDocumentShell` }, t)
			: void 0,
		title: M(n.title, t),
		emptyMessage: M(n.emptyMessage, t),
	});
}
function Lr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Mn, {
		className: n.className,
		title: M(n.title, t),
		description: M(n.description, t),
		accept: n.accept,
		disabled: n.disabled,
		active: n.active,
		onFilesSelected: n.onFilesSelectedAction
			? (e) => {
					Rr(n.onFilesSelectedAction, e, t);
				}
			: void 0,
	});
}
async function Rr(e, t, n) {
	let r = await Promise.all(
		t.map(async (e) => ({
			name: e.name,
			size: e.size,
			type: e.type,
			lastModified: e.lastModified,
			text: await e.text(),
		})),
	);
	oi(
		e,
		{
			componentType: `ContextUploadDropArea`,
			files: r,
			fileNames: r.map((e) => e.name),
			fileCount: r.length,
		},
		n,
	);
}
function zr(e) {
	return typeof e == `object` && !!e && `text` in e;
}
function Br(e) {
	return typeof e == `object` && !!e && `text` in e;
}
function Vr(e, t, n, r) {
	if (t) return (i) => oi(t, { [n]: i, value: i, componentType: e }, r);
}
function Hr(e, t, n) {
	if (t) return (r) => oi(t, { annotationId: r, value: r, componentType: e }, n);
}
function Ur(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(A, {
		className: n.className,
		tone: n.tone,
		transform: n.transform,
		truncate: n.truncate,
		children: Jn(e.children, t),
	});
}
function Wr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Be, { className: n.className, recipe: n.recipe, children: Jn(e.children, t) });
}
function Gr(e, t) {
	let n = e.props ?? {},
		r = n.onRowSelect;
	return (0, d.jsx)(st, {
		className: n.className,
		rows: n.rows,
		columns: n.columns.map((e) => ({
			id: e.id,
			header: M(e.header, t),
			align: e.align,
			maxWidth: e.maxWidth,
			cell: (n) => zn(e.cell, n, (e) => qn(e, t)),
		})),
		getRowKey: (e) => Vn(e, n.getRowKey),
		selectedKey: n.selectedKey == null ? n.selectedKey : String(n.selectedKey),
		emptyMessage: M(n.emptyMessage, t),
		onRowSelect: r
			? (e) => oi(r, { row: e, rowKey: Vn(e, n.getRowKey), componentType: `DataTable` }, t)
			: void 0,
	});
}
function Kr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Je, {
		className: n.className,
		label: M(n.label, t),
		control: qn(n.control, t),
		hint: M(n.hint, t),
	});
}
function qr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Le, {
		className: n.className,
		gap: n.gap,
		justify: n.justify,
		wrap: n.wrap,
		children: Jn(e.children, t),
	});
}
function Jr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(at, {
		className: n.className,
		density: n.density,
		items: n.items.map((e) => ({ key: M(e.key, t), value: M(e.value, t), copyValue: e.copyValue })),
	});
}
function Yr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(ke, {
		className: n.className,
		title: M(n.title, t),
		actions: M(n.actions, t),
		density: n.density,
		fill: n.fill,
		children: Jn(e.children, t),
	});
}
function Xr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(We, {
		className: n.className,
		axis: n.axis,
		style: n.style,
		children: Jn(e.children, t),
	});
}
function Zr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Xe, {
		className: n.className,
		as: n.as,
		label: M(n.label, t),
		caption: M(n.caption, t),
		density: n.density,
		divider: n.divider,
		children: Jn(e.children, t),
	});
}
function Qr(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Qe, {
		className: n.className,
		left: qn(n.left, t),
		right: qn(n.right, t),
		ratio: n.ratio,
		divider: n.divider,
		gutter: n.gutter,
	});
}
function $r(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(rt, {
		className: n.className,
		sidebar: Yn(n.sidebar, t),
		sidebarWidth: n.sidebarWidth,
		contentPadding: n.contentPadding,
		header: Yn(n.header, t),
		footer: Yn(n.footer, t),
		children: Jn(e.children, t),
	});
}
function ei(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(et, {
		className: n.className,
		as: n.as,
		eyebrow: M(n.eyebrow, t),
		counter: M(n.counter, t),
		title: M(n.title, t),
		subtitle: M(n.subtitle, t),
		primary: qn(n.primary, t),
		secondary: Yn(n.secondary, t),
		primarySide: n.primarySide,
		ratio: n.ratio,
		divider: n.divider,
		footer: Yn(n.footer, t),
	});
}
function ti(e, t) {
	let n = e.props ?? {};
	return (0, d.jsxs)(g, {
		className: n.className,
		name: n.name,
		value: n.value,
		disabled: n.disabled,
		children: [
			(n.options ?? []).map((e) =>
				(0, d.jsx)(
					`option`,
					{ value: e.value, disabled: e.disabled, children: M(e.label, t) },
					String(e.value),
				),
			),
			Jn(e.children, t),
		],
	});
}
function ni(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Ne, {
		className: n.className,
		gap: n.gap,
		align: n.align,
		children: Jn(e.children, t),
	});
}
function ri(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Ce, {
		className: n.className,
		status: n.status,
		icon: n.icon,
		children: Jn(e.children, t),
	});
}
function ii(e, t) {
	let n = e.props ?? {};
	return (0, d.jsx)(Ke, {
		items: n.items.map((e) => ({ id: e.id, label: M(e.label, t) })),
		activeId: n.activeId,
		ariaLabel: n.ariaLabel,
		onChange: (e) => {
			n.onChange && oi(n.onChange, { value: e, componentType: `TabList` }, t);
		},
	});
}
function ai(e) {
	let t = e.props ?? {};
	return (0, d.jsx)(v, {
		className: t.className,
		name: t.name,
		value: t.value,
		placeholder: t.placeholder,
		type: t.type,
		disabled: t.disabled,
		min: t.min,
		max: t.max,
		readOnly: !0,
	});
}
function oi(e, t, n) {
	Fn(e, t, n)?.();
}
function si(e, t = {}) {
	let { enabled: n = !0, fetcher: r = fetch } = t,
		[i, a] = (0, O.useState)(null),
		[o, s] = (0, O.useState)(!!(n && e)),
		[c, l] = (0, O.useState)(null),
		[u, d] = (0, O.useState)(0),
		f = (0, O.useCallback)(() => d((e) => e + 1), []);
	return (
		(0, O.useEffect)(() => {
			if (!n || !e) {
				s(!1);
				return;
			}
			let t = new AbortController();
			return (
				s(!0),
				l(null),
				r(e, { signal: t.signal })
					.then(async (e) => {
						if (!e.ok) throw Error(`Widget page request failed: ${e.status} ${e.statusText}`);
						return e.json();
					})
					.then((e) => {
						t.signal.aborted || a(e);
					})
					.catch((e) => {
						t.signal.aborted || (l(e instanceof Error ? e : Error(String(e))), a(null));
					})
					.finally(() => {
						t.signal.aborted || s(!1);
					}),
				() => t.abort()
			);
		}, [n, r, e, u]),
		{ page: i, loading: o, error: c, refresh: f }
	);
}
var ci = [
	{ id: `index`, label: `Overview` },
	{ id: `demo`, label: `Demo` },
	{ id: `actions`, label: `Actions` },
];
function li({ apiBase: e = `/api/widget`, defaultPageId: t = `index` }) {
	let n = bi(t),
		r = e.replace(/\/$/, ``),
		{ page: i, loading: a, error: o, refresh: s } = si(`${r}/pages/${encodeURIComponent(n)}`);
	async function c(e, t) {
		if (e.kind !== `server`) {
			Pn(e, t);
			return;
		}
		let n = await fetch(`${r}/actions/${encodeURIComponent(e.name)}`, {
			method: `POST`,
			headers: { "Content-Type": `application/json` },
			body: JSON.stringify({ payload: e.payload ?? {}, context: t }),
		});
		if (!n.ok) throw Error(`Widget action failed: ${n.status} ${n.statusText}`);
		let i = await n.json();
		i.toast &&
			typeof window < `u` &&
			window.dispatchEvent(new CustomEvent(`widget:toast`, { detail: i })),
			i.refresh && s();
	}
	return a
		? (0, d.jsx)(ke, {
				className: `rag-evaluation-site-state`,
				title: `RAG Evaluation Site`,
				density: `condensed`,
				children: (0, d.jsx)(A, { children: `Loading Widget IR…` }),
			})
		: o
			? (0, d.jsxs)(m, {
					className: `rag-evaluation-site-state`,
					children: [`Failed to load Widget IR page: `, o.message],
				})
			: i
				? ui(i, n, (e, t) => {
						c(e, t);
					})
				: (0, d.jsx)(ke, {
						className: `rag-evaluation-site-state`,
						title: `RAG Evaluation Site`,
						density: `condensed`,
						children: (0, d.jsx)(A, { children: `No Widget IR returned.` }),
					});
}
function ui(e, t, n) {
	let r = pi(e.meta),
		i = (0, d.jsx)(Kn, { node: e.root, onAction: n }),
		a = [
			`rag-evaluation-site-root`,
			di(e, r) ? `rag-evaluation-site-root--shell` : `rag-evaluation-site-root--raw`,
		]
			.filter(Boolean)
			.join(` `);
	if (!di(e, r))
		return (0, d.jsx)(`div`, {
			className: a,
			"data-rag-page": `RagEvaluationSiteApp`,
			"data-page-id": e.id,
			"data-rag-shell": `none`,
			children: i,
		});
	let o = r.navItems?.length ? r.navItems : ci,
		s = r.activeNavItemId ?? t;
	return (0, d.jsx)(`div`, {
		className: a,
		"data-rag-page": `RagEvaluationSiteApp`,
		"data-page-id": e.id,
		"data-rag-shell": `default`,
		children: (0, d.jsx)(He, {
			className: `rag-evaluation-site-shell`,
			header: (0, d.jsx)(lt, {
				brand: (0, d.jsx)(`span`, {
					className: `rag-evaluation-site-brand`,
					children: e.title || `RAG Evaluation Site`,
				}),
				items: o.map((e) => ({ id: e.id, label: _i(e.label) })),
				activeItemId: s,
				onItemSelect: (e) => {
					let t = o.find((t) => t.id === e);
					if (t?.action) {
						n(t.action, { value: e, componentType: `AppNav` });
						return;
					}
					yi(e);
				},
			}),
			children: (0, d.jsx)(`div`, {
				className: vi(r.maxWidth),
				"data-rag-layout": `PageContent`,
				children: i,
			}),
		}),
	});
}
function di(e, t) {
	return !(t.shell === `none` || fi(e.root));
}
function fi(e) {
	return e.kind === `component` && e.type === `AppShell`;
}
function pi(e) {
	return e
		? {
				shell: mi(e.shell) ? e.shell : void 0,
				activeNavItemId: typeof e.activeNavItemId == `string` ? e.activeNavItemId : void 0,
				navItems: Array.isArray(e.navItems) ? e.navItems.filter(gi) : void 0,
				maxWidth: hi(e.maxWidth) ? e.maxWidth : void 0,
			}
		: {};
}
function mi(e) {
	return e === `auto` || e === `none` || e === `app`;
}
function hi(e) {
	return e === `none` || e === `content` || e === `wide`;
}
function gi(e) {
	if (!e || typeof e != `object`) return !1;
	let t = e;
	return typeof t.id == `string` && t.label !== void 0;
}
function _i(e) {
	return e && typeof e == `object` && `kind` in e ? (0, d.jsx)(Kn, { node: e }) : String(e ?? ``);
}
function vi(e) {
	return [`rag-evaluation-site-content`, `rag-evaluation-site-content--${e ?? `wide`}`].join(` `);
}
function yi(e) {
	if (typeof window > `u`) return;
	let t = new URL(window.location.href);
	(t.pathname = `/pages/${encodeURIComponent(e)}`),
		t.searchParams.delete(`page`),
		window.location.assign(t.toString());
}
function bi(e) {
	if (typeof window > `u`) return e;
	let t = new URL(window.location.href),
		n = t.searchParams.get(`page`);
	if (n) return n;
	let r = t.pathname.split(`/`).filter(Boolean);
	return r[0] === `pages` && r[1] ? r[1] : e;
}
var xi = document.getElementById(`root`);
if (!xi) throw Error(`RagEvaluationSiteApp requires an element with id="root"`);
(0, Nn.createRoot)(xi).render((0, d.jsx)(O.StrictMode, { children: (0, d.jsx)(li, {}) }));

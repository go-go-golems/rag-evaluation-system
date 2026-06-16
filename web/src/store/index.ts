import { configureStore } from "@reduxjs/toolkit";
import { ragApi } from "../services/api";
import { useDispatch, useSelector, TypedUseSelectorHook } from "react-redux";

export const store = configureStore({
	reducer: {
		[ragApi.reducerPath]: ragApi.reducer,
	},
	middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(ragApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export const typedHooks = {
	useDispatch: useDispatch<AppDispatch>,
	useSelector: useSelector as TypedUseSelectorHook<RootState>,
};

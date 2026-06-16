import { useDispatch, useSelector, TypedUseSelectorHook } from "react-redux";
import type { RootState, AppDispatch } from "../store";

export const useAppDispatch = useDispatch<AppDispatch>;
export const useAppSelector = useSelector as TypedUseSelectorHook<RootState>;

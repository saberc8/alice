import { lazy, Suspense } from "react";
import { LineLoading } from "@/components/loading";

const Pages = import.meta.glob("/src/pages/**/*.tsx");
const lazyComponentCache = new Map<string, React.LazyExoticComponent<any>>();

export const loadComponentFromPath = (path: string) => {
	const pathArr = path.split("/");
	pathArr.unshift("/src");

	if (!pathArr.includes(".tsx")) {
		return pathArr.push("index.tsx");
	}
	return Pages[pathArr.join("/")];
};

export const Component = (path = "", props?: any): React.ReactNode => {
	if (!path) {
		console.warn("Component path is empty");
		return null;
	}

	// 确保路径以 /src 开始
	const normalizedPath = path.startsWith('/src') ? path : `/src${path}`;
	
	// 尝试两种文件模式：直接 .tsx 文件和 index.tsx
	let importFn = Pages[`${normalizedPath}.tsx`];
	if (!importFn) {
		importFn = Pages[`${normalizedPath}/index.tsx`];
	}
	
	if (!importFn) {
		console.error(`组件未找到: ${path}`);
		console.log(`尝试的路径:`, [`${normalizedPath}.tsx`, `${normalizedPath}/index.tsx`]);
		
		// 显示可能的匹配项
		const pathSegments = path.split('/').filter(Boolean);
		const lastSegment = pathSegments[pathSegments.length - 1];
		const possibleMatches = Object.keys(Pages).filter(key => 
			key.toLowerCase().includes(lastSegment?.toLowerCase() || '')
		);
		
		if (possibleMatches.length > 0) {
			console.log(`可能的匹配:`, possibleMatches);
		}
		
		return (
			<div className="p-4 text-center">
				<h3 className="text-lg font-semibold text-red-600 mb-2">组件加载失败</h3>
				<p className="text-gray-600">找不到组件: <code className="bg-gray-100 px-2 py-1 rounded">{path}</code></p>
				<p className="text-sm text-gray-500 mt-2">请检查组件路径是否正确</p>
			</div>
		);
	}

	let Element = lazyComponentCache.get(path);
	if (!Element) {
		Element = lazy(importFn as any);
		lazyComponentCache.set(path, Element);
	}
	
	return (
		<Suspense fallback={<LineLoading />}>
			<Element {...props} />
		</Suspense>
	);
};

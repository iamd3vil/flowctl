/**
 * @typedef {{ label: string, url?: string }} BreadcrumbItem
 */

/**
 * @param {string} namespace
 * @param {string} prefix
 * @returns {BreadcrumbItem | null}
 */
function buildGroupCrumb(namespace, prefix) {
    if (!prefix) {
        return null;
    }

    return {
        label: prefix,
        url: `/view/${namespace}/flows?group=${encodeURIComponent(prefix)}`,
    };
}

/**
 * @param {string} namespace
 * @param {string} flowName
 * @param {string} prefix
 * @param {{ flowUrl?: string, trailingLabel?: string }=} options
 * @returns {BreadcrumbItem[]}
 */
export function buildFlowBreadcrumbs(namespace, flowName, prefix, options = {}) {
    const breadcrumbs = [
        { label: namespace },
        { label: "Flows", url: `/view/${namespace}/flows` },
    ];

    const groupCrumb = buildGroupCrumb(namespace, prefix);
    if (groupCrumb) {
        breadcrumbs.push(groupCrumb);
    }

    breadcrumbs.push(
        options.flowUrl
            ? { label: flowName, url: options.flowUrl }
            : { label: flowName },
    );

    if (options.trailingLabel) {
        breadcrumbs.push({ label: options.trailingLabel });
    }

    return breadcrumbs;
}

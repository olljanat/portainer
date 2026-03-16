import { useStore } from 'zustand';

import { useEnvironmentId } from '@/react/hooks/useEnvironmentId';
import { useNamespaces } from '@/react/kubernetes/namespaces/queries';
import { useAuthorizations } from '@/react/hooks/useUser';
import Route from '@/assets/ico/route.svg?c';

import { Datatable } from '@@/datatables';
import { createPersistedStore } from '@@/datatables/types';
import { useSearchBarState } from '@@/datatables/SearchBar';

import { useIngresses } from '../queries';

import { useColumns } from './columns';

import '../style.css';

const storageKey = 'ingressClassesNameSpace';

const settingsStore = createPersistedStore(storageKey);

export function IngressDatatable() {
  const environmentId = useEnvironmentId();

  const nsResult = useNamespaces(environmentId);
  const ingressesQuery = useIngresses(
    environmentId,
    Object.keys(nsResult?.data || {})
  );

  const columns = useColumns();
  const settings = useStore(settingsStore);
  const [search, setSearch] = useSearchBarState(storageKey);

  return (
    <Datatable
      dataset={ingressesQuery.data || []}
      columns={columns}
      isLoading={ingressesQuery.isLoading}
      emptyContentLabel="No supported ingresses found"
      title="Ingresses"
      titleIcon={Route}
      getRowId={(row) => row.Name + row.Type + row.Namespace}
      renderTableActions={tableActions}
      disableSelect={useCheckboxes()}
      initialPageSize={settings.pageSize}
      onPageSizeChange={settings.setPageSize}
      initialSortBy={settings.sortBy}
      onSortByChange={settings.setSortBy}
      searchValue={search}
      onSearchChange={setSearch}
    />
  );

  function tableActions() {
    return <div className="ingressDatatable-actions" />;
  }

  function useCheckboxes() {
    return !useAuthorizations(['K8sIngressesW']);
  }
}

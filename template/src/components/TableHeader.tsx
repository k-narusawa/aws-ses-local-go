export const TableHeader: React.FC = () => {
  return (
    <>
      <thead className="bg-gray-50">
        <tr>
          <th scope="col" className="py-3 px-4 pe-0">
            <div className="flex items-center h-5">
              <input
                id="hs-table-pagination-checkbox-all"
                type="checkbox"
                className="border-gray-200 rounded text-blue-600 focus:ring-blue-500"
              />
            </div>
          </th>
          <th
            scope="col"
            className="px-6 py-3 text-start text-xs font-medium text-gray-500 uppercase"
          >
            Subject
          </th>
          <th
            scope="col"
            className="px-6 py-3 text-start text-xs font-medium text-gray-500 uppercase"
          >
            From
          </th>
          <th
            scope="col"
            className="px-6 py-3 text-start text-xs font-medium text-gray-500 uppercase"
          >
            To
          </th>
          <th
            scope="col"
            className="px-6 py-3 text-end text-xs font-medium text-gray-500 uppercase"
          >
            Action
          </th>
        </tr>
      </thead>
    </>
  );
};

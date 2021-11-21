import { useNavigate } from "react-router";

function FoldersTable(props) {
  let navigate = useNavigate();
  return (
    <table>
      <thead>
        <tr>
          <th>Folder name</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        {props.folders.map((folder) => (
          <tr key={folder.id}>
            <td
              onClick={() => navigate(`/search-folders/${folder.id}`)}
              className="clickable"
            >
              {folder.foldername}
            </td>
            <td onClick={() => navigate(`/search-folders/${folder.id}`)}>
              <button>View</button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default FoldersTable;

import { useNavigate } from "react-router";
import { Tr } from "../translate";

function FoldersTable(props) {
  let navigate = useNavigate();
  if (props.folders.length === 0) {
    return null;
  }
  return (
    <table>
      <thead>
        <tr>
          <th>{Tr("Folder name")}</th>
          <th>{Tr("Action")}</th>
        </tr>
      </thead>
      <tbody>
        {props.folders.map((folder) => (
          <tr key={folder.id}>
            <td
              onClick={() => navigate(`/folders/${folder.id}`)}
              className="clickable"
            >
              {folder.foldername}
            </td>
            <td onClick={() => navigate(`/folders/${folder.id}`)}>
              <button>{Tr("View")}</button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default FoldersTable;

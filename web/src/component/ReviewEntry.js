import { Link } from "react-router-dom";
import { useNavigate } from "react-router";
import { convertIntToDateTime } from "./Common";

function ReviewEntry(props) {
  let navigate = useNavigate();
  return (
    <div>
      <h4>
        <Link to={`/manage/users/${props.review.user.id}`}>
          @{props.review.user.username}
        </Link>{" "}
        wrote on {convertIntToDateTime(props.review.created_at)}{" "}
      </h4>
      <p>{props.review.content}</p>
      {(props.user.role === 1 || props.review.user.id === props.user.id) &&
        props.user.role != 0 && (
          <button
            onClick={() => {
              navigate(`/manage/reviews/${props.review.id}`);
            }}
          >
            Edit
          </button>
        )}
    </div>
  );
}

export default ReviewEntry;

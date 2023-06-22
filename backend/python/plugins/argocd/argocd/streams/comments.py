import pydevlake as dl
import pydevlake.domain_layer.ticket as ticket

from typing import Iterable, Tuple

class Comments(dl.Substream):
    tool_model = ticket.IssueComment
    domain_models = [ticket.IssueComment]
    parent_stream = ticket.Issues

    def collect(self, state, context, parent: ticket.Issue) -> Iterable[Tuple[object, dict]]:
        ...

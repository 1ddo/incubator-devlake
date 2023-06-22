from typing import Iterable, Tuple
from pydevlake import Stream, DomainType
from pydevlake.model import ToolModel
from argocd.models import User, ArgocdAPI

import pydevlake.domain_layer.crossdomain as cross


class Users(Stream):
    tool_model = ToolUser
    domain_models = [cross.User]

    def collect(self, state, context) -> Iterable[Tuple[object, dict]]:
        api = ArgocdAPI(context.connection.url)
        for user in api.users().json():
            yield user, state

    def extract(self, raw_data) -> ToolUser:
        return ToolUser(
            id=raw_data["ID"],
            name=raw_data["Name"],
            email=raw_data["Email"]
        )

    def convert(self, user: ToolUser, context) -> Iterable[DomainUser]:
        yield DomainUser(
            id=user.id,
            name=user.name,
            email=user.email,
        )
    
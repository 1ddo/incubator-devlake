#!/usr/bin/env python

from typing import Iterable
from argocd.api import ArgocdAPI
from pydevlake import Plugin, RemoteScopeGroup, DomainType, ScopeConfigPair
from pydevlake.domain_layer.code import Repo
from pydevlake.domain_layer.devops import CicdScope
from pydevlake.pipeline_tasks import gitextractor, refdiff
from pydevlake.api import APIException

import pydevlake as dl

class ArgocdPluginConnection(dl.Connection):
    """
    connection that groups the parameters that your plugin needs to collect data, e.g. 
    the url and credentials to connect to the datasource

    Yields:
        _type_: _description_
    """

    url: str
    org: str
    token: str


#class ArgocdPluginTransformationRule(dl.TransformationRule):
#    issue_type_regex: str


class ArgocdPluginToolScope(dl.ToolScope):
    """
    tool layer scope type that represents the top-level entity of this plugin, e.g. 
    a board, a repository, a project, etc.

    Yields:
        _type_: _description_
    """

    id: str
    name: str
    description: str
    url: str

class ArgocdPlugin(dl.Plugin):
    """_summary_

    Args:
        dl (_type_): _description_

    Yields:
        _type_: _description_
    """
    connection_type = ArgocdPluginConnection
    connection_type.url = 'https://datausa.io/api/data?drilldowns=Nation&measures=Population'
    #transformation_rule_type =  ArgocdPluginTransformationRule
    tool_scope_type = ArgocdPluginToolScope
    streams = []

    def domain_scopes(self, tool_scope: ArgocdPluginToolScope) -> Iterable[dl.DomainScope]:
        yield CicdScope(
            name=tool_scope.name,
            description=tool_scope.description,
            url=tool_scope.url,
        )

    def remote_scope_groups(self, connection: ArgocdPluginConnection) -> Iterable[dl.RemoteScopeGroup]:
        api = ArgocdAPI(connection)
        apps = api.applications(connection.org)

        for app in apps:
            items = app['items']

            for item in items:
                metadata = item['metadata']
                app_name = metadata['name']
                #status = metadata['status']
                #histories = status['history']

                #for history in histories:
                    #rev = history['revision']
                    #depAt = history['deployAt']
                    #depStartAt = history['deployStartedAt']
                    #src = history['source']
                    #repoURL = src['repoURL']
                    #path = src['path']
                    #targetRev  = src['targetRevision']

                yield RemoteScopeGroup(
                    id=f'{connection.org}/{app_name}',
                    name=app_name
                )

    def remote_scopes(self, connection, group_id: str) -> Iterable[ArgocdPluginToolScope]:
        api = ArgocdAPI(connection)
        apps = api.applications(connection.org)

        for app in apps:
            items = app['items']

            for item in items:
                metadata = item['metadata']
                app_name = metadata['name']
                uid = metadata['uid']
                #status = metadata['status']
                #histories = status['history']

                yield ArgocdPluginToolScope(
                    id=uid,
                    name=app_name,
                    #description=raw_scope['description'],
                    #url=raw_scope['url'],
                )

    def test_connection(self, connection: ArgocdPluginConnection):
        api = ArgocdAPI(connection)
        apps = api.applications(connection.org)
        response = apps.test_connection

        if response.status != 401:
            raise Exception("Invalid credentials")
        if response.status != 200:
            raise Exception(f"Connection error {response}") 

if __name__ == '__main__':
    ArgocdPlugin.start()

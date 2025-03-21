/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import { Switch, Route, Redirect, Router } from 'react-router-dom';
import { ErrorLayout, BaseLayout } from '@/layouts';
import {
  LoginPage,
  OfflinePage,
  DBMigratePage,
  ConnectionHomePage,
  ConnectionDetailPage,
  ProjectHomePage,
  ProjectDetailPage,
  BlueprintHomePage,
  BlueprintDetailPage,
  BlueprintConnectionDetailPage,
} from '@/pages';
import { history } from '@/utils';

function App() {
  return (
    <Router history={history}>
      <Switch>
        <Route exact path="/login" component={() => <LoginPage />} />

        <Route
          exact
          path="/offline"
          component={() => (
            <ErrorLayout>
              <OfflinePage />
            </ErrorLayout>
          )}
        />

        <Route
          exact
          path="/db-migrate"
          component={() => (
            <ErrorLayout>
              <DBMigratePage />
            </ErrorLayout>
          )}
        />

        <Route
          path="/"
          component={() => (
            <BaseLayout>
              <Switch>
                <Route exact path="/" component={() => <Redirect to="/connections" />} />
                <Route exact path="/connections" component={() => <ConnectionHomePage />} />
                <Route exact path="/connections/:plugin/:id" component={() => <ConnectionDetailPage />} />
                <Route exact path="/projects" component={() => <ProjectHomePage />} />
                <Route exact path="/projects/:pname" component={() => <ProjectDetailPage />} />
                <Route exact path="/projects/:pname/:unique" component={() => <BlueprintConnectionDetailPage />} />
                <Route exact path="/blueprints" component={() => <BlueprintHomePage />} />
                <Route exact path="/blueprints/:id" component={() => <BlueprintDetailPage />} />
                <Route exact path="/blueprints/:bid/:unique" component={() => <BlueprintConnectionDetailPage />} />
              </Switch>
            </BaseLayout>
          )}
        />
      </Switch>
    </Router>
  );
}

export default App;

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

import { useState, useEffect } from 'react';
import { Tag, Intent, Switch, InputGroup } from '@blueprintjs/core';

import * as S from './styled';
import { ExternalLink, HelpTooltip } from '@/components';

interface Props {
  entities: string[];
  transformation: any;
  setTransformation: React.Dispatch<React.SetStateAction<any>>;
}

export const JenkinsTransformation = ({ entities, transformation, setTransformation }: Props) => {
  const [enableCICD, setEnableCICD] = useState(true);

  useEffect(() => {
    if (!transformation.deploymentPattern) {
      setEnableCICD(false);
    }
  }, [transformation]);

  const handleChangeCICDEnable = (e: React.FormEvent<HTMLInputElement>) => {
    const checked = (e.target as HTMLInputElement).checked;

    if (!checked) {
      setTransformation({
        ...transformation,
        deploymentPattern: undefined,
        productionPattern: undefined,
      });
    }

    setEnableCICD(checked);
  };

  return (
    <S.Transformation>
      {entities.includes('CICD') && (
        <S.CICD>
          <h2>CI/CD</h2>
          <h3>
            <span>Deployment</span>
            <Tag minimal intent={Intent.PRIMARY} style={{ marginLeft: 8 }}>
              DORA
            </Tag>
            <div className="switch">
              <span>Enable</span>
              <Switch alignIndicator="right" inline checked={enableCICD} onChange={handleChangeCICDEnable} />
            </div>
          </h3>
          {enableCICD && (
            <>
              <p>
                Use Regular Expression to define Deployments in DevLake in order to measure DORA metrics.{' '}
                <ExternalLink link="https://devlake.apache.org/docs/Configuration/GitHub#step-3---adding-transformation-rules-optional">
                  Learn more
                </ExternalLink>
              </p>
              <div style={{ marginTop: 16 }}>Convert a Jenkins Build as a DevLake Deployment when: </div>
              <div className="text">
                <span>
                  The name of the <strong>Jenkins job</strong> or <strong>one of its stages</strong> matches
                </span>
                <InputGroup
                  style={{ width: 200, margin: '0 8px' }}
                  placeholder="(deploy|push-image)"
                  value={transformation.deploymentPattern ?? ''}
                  onChange={(e) =>
                    setTransformation({
                      ...transformation,
                      deploymentPattern: e.target.value,
                      productionPattern: !e.target.value ? '' : transformation.productionPattern,
                    })
                  }
                />
                <i style={{ color: '#E34040' }}>*</i>
                <HelpTooltip content="Jenkins Builds: https://www.jenkins.io/doc/pipeline/steps/pipeline-build-step/" />
              </div>
              <div className="text">
                <span>If the name also matches</span>
                <InputGroup
                  style={{ width: 120, margin: '0 8px' }}
                  disabled={!transformation.deploymentPattern}
                  placeholder="prod(.*)"
                  value={transformation.productionPattern ?? ''}
                  onChange={(e) =>
                    setTransformation({
                      ...transformation,
                      productionPattern: e.target.value,
                    })
                  }
                />
                <span>, this Deployment is a ‘Production Deployment’</span>
                <HelpTooltip content="If you leave this field empty, all DevLake Deployments will be tagged as in the Production environment. " />
              </div>
            </>
          )}
        </S.CICD>
      )}
    </S.Transformation>
  );
};

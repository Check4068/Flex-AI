FROM openeuler/openeuler:22.03-lts-sp2

COPY vc-scheduler /vc-scheduler
COPY huawei-*.so plugins/

RUN chmod -R root:root plugins vc-scheduler && \
    chmod 500 plugins vc-scheduler && \
    chmod 400 plugins/*.so
USER root